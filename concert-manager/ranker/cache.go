package ranker

import (
	"concert-manager/domain"
	"concert-manager/external"
	"concert-manager/file"
	"concert-manager/log"
	"fmt"
	"strings"
	"sync"
	"time"
)

type spotifyService interface {
	SavedTracks() ([]external.Track, error)
	TopTracks(external.TimeRange) ([]external.Track, error)
	TopArtists(external.TimeRange) ([]external.Artist, error)
}

type artistProvider interface {
	SimilarArtists(string) ([]external.RankedArtist, error)
}

type ArtistRankCache struct {
	MusicSvc       spotifyService
	ArtistProvider artistProvider
	ranks          map[string]domain.ArtistRank
	lastRefresh    time.Time
	refreshing     bool
	refreshMutex   sync.Mutex
	calculator     *RankCalculator
}

type RankCacheFile struct {
	Timestamp time.Time                    `json:"timestamp"`
	Version   string                       `json:"version"`
	Ranks     map[string]domain.ArtistRank `json:"ranks"`
}

const (
	rankCacheVersion = "1.0"
	rankCacheFile    = "ranks.json"
)

var rankTTL, _ = time.ParseDuration("168h")

func (c *ArtistRankCache) Rank(artist domain.Artist) domain.ArtistRank {
	if c.lastRefresh.IsZero() {
		c.DoRefresh()
	} else if time.Since(c.lastRefresh) > rankTTL {
		go c.DoRefresh()
	}
	return c.ranks[toKey(artist.Name)]
}

func (c *ArtistRankCache) DoRefresh() {
	c.refreshMutex.Lock()
	if c.refreshing {
		c.refreshMutex.Unlock()
		return
	}
	c.refreshing = true
	c.refreshMutex.Unlock()

	log.Info("Refreshing artist ranks")

	if c.calculator == nil {
		c.calculator = &RankCalculator{
			MusicSvc:       c.MusicSvc,
			ArtistProvider: c.ArtistProvider,
		}
	}

	newRanks, err := c.calculator.CalculateRanks()
	if err != nil {
		log.Alert("Failed to refresh artist ranks", err)
	} else {
		c.ranks = newRanks
		log.Info("Successfully refreshed artist ranks")
		c.saveRanksToFile()
	}

	c.lastRefresh = time.Now().Round(0)
	c.refreshMutex.Lock()
	c.refreshing = false
	c.refreshMutex.Unlock()
}

func (c *ArtistRankCache) InitializeFromFile() error {
	filePath, err := file.GetCacheFilePath(rankCacheFile)
	if err != nil {
		return fmt.Errorf("failed to get cache file path: %w", err)
	}

	if !file.FileExists(filePath) {
		log.Debug("Rank cache file does not exist")
		c.DoRefresh()
		return nil
	}

	log.Info("Loading artist ranks from cache file")
	err = c.loadRanksFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load ranks from file: %v", err)
	}

	if file.IsFileStale(filePath, rankTTL) {
		log.Info("Rank cache file is stale, starting background refresh")
		go c.DoRefresh()
	}

	return nil
}

func (c *ArtistRankCache) initializeEmpty() {
	c.ranks = make(map[string]domain.ArtistRank)
	c.lastRefresh = time.Time{}
}

func (c *ArtistRankCache) loadRanksFromFile(filePath string) error {
	var cacheFile RankCacheFile
	err := file.ReadJSONFile(filePath, &cacheFile)
	if err != nil {
		return err
	}

	if cacheFile.Version != rankCacheVersion {
		go c.DoRefresh()
		return nil
	}

	c.ranks = cacheFile.Ranks
	if c.ranks == nil {
		c.ranks = make(map[string]domain.ArtistRank)
	}
	c.lastRefresh = cacheFile.Timestamp

	log.Infof("Loaded %d artist ranks from cache file", len(c.ranks))
	return nil
}

func (c *ArtistRankCache) saveRanksToFile() {
	filePath, err := file.GetCacheFilePath(rankCacheFile)
	if err != nil {
		log.Errorf("Failed to get cache file path for saving: %v", err)
		return
	}

	cacheFile := RankCacheFile{
		Timestamp: time.Now().Round(0),
		Version:   rankCacheVersion,
		Ranks:     c.ranks,
	}

	err = file.WriteJSONFile(filePath, cacheFile)
	if err != nil {
		log.Errorf("Failed to save ranks to file: %v", err)
		return
	}

	log.Infof("Successfully saved %d artist ranks to cache file", len(c.ranks))
}

func toKey(name string) string {
	return strings.ToLower(name)
}

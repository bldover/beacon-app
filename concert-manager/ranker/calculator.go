package ranker

import (
	"concert-manager/domain"
	"concert-manager/external"
	"concert-manager/log"
	"slices"
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
}

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
	err := c.refreshRanks()
	if err != nil {
		// this would only happen for catastrophic failures or large Spotify rate limits, so don't retry
		log.Alert("Failed to refresh artist ranks", err)
	} else {
		log.Info("Successfully refreshed artist ranks")
	}

	c.lastRefresh = time.Now().Round(0)
	c.refreshMutex.Lock()
	c.refreshing = false
	c.refreshMutex.Unlock()
}

const (
	trackRankCeilingPercent  = 0.9
	savedTrackFactor         = 10.
	topTrackLongTermFactor   = 12.
	topTrackMediumTermFactor = 10.
	topTrackShortTermFactor  = 5.

	topArtistLongTermFactor   = 12.
	topArtistMediumTermFactor = 10.
	topArtistShortTermFactor  = 5.

	featuredArtistFactor = 0.3
	similarArtistFactor  = 0.15
)

func (c *ArtistRankCache) refreshRanks() error {
	rawSavedTracks, err := c.MusicSvc.SavedTracks()
	if err != nil {
		return err
	}
	rawTopTracksLongTerm, err := c.MusicSvc.TopTracks(external.LongTerm)
	if err != nil {
		return err
	}
	rawTopTracksMediumTerm, err := c.MusicSvc.TopTracks(external.MediumTerm)
	if err != nil {
		return err
	}
	rawTopTracksShortTerm, err := c.MusicSvc.TopTracks(external.ShortTerm)
	if err != nil {
		return err
	}
	rawTopArtistsLongTerm, err := c.MusicSvc.TopArtists(external.LongTerm)
	if err != nil {
		return err
	}
	rawTopArtistsMediumTerm, err := c.MusicSvc.TopArtists(external.MediumTerm)
	if err != nil {
		return err
	}
	rawTopArtistsShortTerm, err := c.MusicSvc.TopArtists(external.ShortTerm)
	if err != nil {
		return err
	}

	savedTracks := mapTracks(rawSavedTracks)
	topTracksLongTerm := mapTracks(rawTopTracksLongTerm)
	topTracksMediumTerm := mapTracks(rawTopTracksMediumTerm)
	topTracksShortTerm := mapTracks(rawTopTracksShortTerm)

	topArtistsLongTerm := mapArtists(rawTopArtistsLongTerm)
	topArtistsMediumTerm := mapArtists(rawTopArtistsMediumTerm)
	topArtistsShortTerm := mapArtists(rawTopArtistsShortTerm)

	prevRanks := c.ranks
	c.ranks = make(map[string]domain.ArtistRank, len(topArtistsLongTerm)*15)

	c.normalizeTrackRanks(topTracksLongTerm)
	c.normalizeTrackRanks(topTracksMediumTerm)
	c.normalizeTrackRanks(topTracksShortTerm)
	c.updateRankForSavedTracks(savedTracks, savedTrackFactor, trackRankCeilingPercent)
	c.updateRankForRankedTracks(topTracksLongTerm, topTrackLongTermFactor, trackRankCeilingPercent)
	c.updateRankForRankedTracks(topTracksMediumTerm, topTrackMediumTermFactor, trackRankCeilingPercent)
	c.updateRankForRankedTracks(topTracksShortTerm, topTrackShortTermFactor, trackRankCeilingPercent)

	c.normalizeArtistRanks(topArtistsLongTerm)
	c.normalizeArtistRanks(topArtistsMediumTerm)
	c.normalizeArtistRanks(topArtistsShortTerm)
	c.updateRankForArtists(topArtistsLongTerm, topArtistLongTermFactor)
	c.updateRankForArtists(topArtistsMediumTerm, topArtistMediumTermFactor)
	c.updateRankForArtists(topArtistsShortTerm, topArtistShortTermFactor)

	err = c.populateSimilarArtistRanks(topArtistsLongTerm)
	if err != nil {
		c.ranks = prevRanks
		return err
	}

	c.normalizeRanks()
	c.logRanks()

	return nil
}

// from spotify client, tracks ordered from highest to lowest rank
// convert these to (0, 1], where top tracks have rank 1
func (c *ArtistRankCache) normalizeTrackRanks(topTracks []rankedTrack) {
	max := float64(len(topTracks))
	for i := range topTracks {
		track := &topTracks[i]
		track.rank = (max - float64(i)) / max
	}
}

// from spotify client, artists ordered from highest to lowest rank
// convert these to (0, 1], where top artists have rank 1
func (c *ArtistRankCache) normalizeArtistRanks(topArtists []rankedArtist) {
	max := float64(len(topArtists))
	for i := range topArtists {
		artist := &topArtists[i]
		artist.rank = (max - float64(i)) / max
	}
}

func (c *ArtistRankCache) normalizeRanks() {
	maxRank := 0.
	for _, rankData := range c.ranks {
		maxRank = max(maxRank, rankData.Rank)
	}
	for artist, rankData := range c.ranks {
		rankData.Rank /= maxRank
		c.ranks[artist] = rankData
	}
}

func (c *ArtistRankCache) updateRankForSavedTracks(tracks []rankedTrack, factor float64, rankCeilPerc float64) {
	tempRanks := make(map[string]float64, len(tracks)*2)
	maxRank := 0.
	for _, track := range tracks {
		artist := track.Artists[0]
		key := toKey(artist.Name)
		tempRanks[key] += 1
		maxRank = max(tempRanks[key], maxRank)
	}

	rankCeiling := maxRank * rankCeilPerc
	for artist, rank := range tempRanks {
		adjRank := min(rank, rankCeiling)
		ceiledRank := (adjRank / rankCeiling) * factor

		key := toKey(artist)
		rankData := c.ranks[key]
		rankData.Rank += ceiledRank
		c.ranks[key] = rankData
	}
}

func (c *ArtistRankCache) updateRankForRankedTracks(tracks []rankedTrack, factor float64, rankCeilPerc float64) {
	tempRanks := make(map[string]float64, len(tracks)*2)
	maxRank := 0.
	for _, track := range tracks {
		rank := track.rank
		for i, artist := range track.Artists {
			if i > 0 {
				rank *= featuredArtistFactor
			}
			key := toKey(artist.Name)
			tempRanks[key] += rank
			maxRank = max(tempRanks[key], maxRank)
		}
	}

	rankCeiling := maxRank * rankCeilPerc
	for artist, rank := range tempRanks {
		adjRank := min(rank, rankCeiling)
		normalizedRank := (adjRank / rankCeiling) * factor

		key := toKey(artist)
		rankData := c.ranks[key]
		rankData.Rank += normalizedRank
		c.ranks[key] = rankData
	}
}

func (c *ArtistRankCache) updateRankForArtists(artists []rankedArtist, weight float64) {
	for _, artist := range artists {
		key := toKey(artist.Name)
		rankData := c.ranks[key]
		rankData.Rank += artist.rank * weight
		c.ranks[key] = rankData
	}
}

func (c *ArtistRankCache) populateSimilarArtistRanks(artists []rankedArtist) error {
	log.Infof("Retrieving similar artist data for %v artists\n", len(artists))
	// need this so we don't reference increasing ranks of known artists as we iterate
	similarArtistRanks := make(map[string]float64, len(artists)*5)

	for _, knownArtist := range artists {
		similarArtists, err := c.ArtistProvider.SimilarArtists(knownArtist.Name)
		if err != nil {
			log.Errorf("Failed to find similar artists for %v, %v", knownArtist, err)
			return err
		}
		for _, similarArtist := range similarArtists {
			if similarArtist.Name == "" {
				continue
			}
			knownArtistData := c.ranks[toKey(knownArtist.Name)]
			calcRankInc := knownArtistData.Rank * similarArtist.Rank * similarArtistFactor

			key := toKey(similarArtist.Name)
			similarArtistRanks[key] += calcRankInc
			similarRankData := c.ranks[key]
			if similarRankData.Related == nil {
				similarRankData.Related = []string{}
			}
			similarRankData.Related = append(similarRankData.Related, knownArtist.Name)
			c.ranks[key] = similarRankData
		}
	}

	for name, rank := range similarArtistRanks {
		key := toKey(name)
		rankData := c.ranks[key]
		rankData.Rank += rank
		c.ranks[key] = rankData
	}
	log.Info("Finished retrieving similar artist data")
	return nil
}

func toKey(name string) string {
	return strings.ToLower(name)
}

type rankedTrack struct {
	external.Track
	rank float64
}

func mapTracks(tracks []external.Track) []rankedTrack {
	ranked := []rankedTrack{}
	for _, track := range tracks {
		ranked = append(ranked, rankedTrack{track, 0})
	}
	return ranked
}

type rankedArtist struct {
	external.Artist
	rank float64
}

func mapArtists(artists []external.Artist) []rankedArtist {
	ranked := []rankedArtist{}
	for _, artist := range artists {
		ranked = append(ranked, rankedArtist{artist, 0})
	}
	return ranked
}

func (c *ArtistRankCache) logRanks() {
	if !log.IsDebug() {
		return
	}

	log.Debug("Logging all artist ranks")
	artists := make([]string, len(c.ranks))
	for artist := range c.ranks {
		artists = append(artists, artist)
	}

	slices.SortFunc(artists, func(a string, b string) int {
		aRank := c.ranks[a].Rank
		bRank := c.ranks[b].Rank
		if aRank < bRank {
			return -1
		} else if aRank > bRank {

			return 1
		} else {
			return 0
		}
	})

	for _, artist := range artists {
		log.Debugf("%s: %v", artist, c.ranks[artist])
	}
}

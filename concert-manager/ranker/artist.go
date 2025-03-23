package ranker

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/spotify"
	"slices"
	"strings"
	"sync"
	"time"
)

type ArtistRanker struct {
	MusicSvc musicService
	ranks map[string]rankData
	lastRefresh time.Time
	refreshing bool
	refreshMutex sync.Mutex
}

type musicService interface {
    GetSavedTracks() ([]spotify.Track, error)
	GetTopTracks(spotify.TimeRange) ([]spotify.RankedTrack, error)
	GetTopArtists(spotify.TimeRange) ([]spotify.RankedArtist, error)
	GetRelatedArtistsBatch([]spotify.Artist) (map[spotify.Artist][]spotify.Artist, error)
}

type rankData struct {
    rank float64
	related []string
}

const (
	trackRankCeilingPercent = 0.9
	savedTrackFactor = 10.
	topTrackLongTermFactor = 12.
	topTrackMediumTermFactor = 10.
	topTrackShortTermFactor = 5.

	topArtistLongTermFactor = 12.
	topArtistMediumTermFactor = 10.
	topArtistShortTermFactor = 5.

	featuredArtistFactor = 0.3
	relatedArtistFactor = 0.15
)

var rankTTL, _ = time.ParseDuration("168h")

func (r *ArtistRanker) Rank(artist data.Artist) data.ArtistRank {
	artistRank := data.ArtistRank{Artist: artist, Related: []string{}}
	if artist.Name == "" {
		return artistRank
	}

	if r.lastRefresh.IsZero() {
		r.DoRefresh()
	} else if time.Since(r.lastRefresh) > rankTTL {
		go r.DoRefresh()
	}

	if rankData, ok := r.ranks[toKey(artist.Name)]; ok {
		artistRank.Rank = rankData.rank
		artistRank.Related = append(artistRank.Related, rankData.related...)
		return artistRank
	}

	log.Debugf("Ranked artist %v\n", artistRank)
 	return artistRank
}

func (r *ArtistRanker) DoRefresh() {
	r.refreshMutex.Lock()
	if r.refreshing {
		r.refreshMutex.Unlock()
		return
	}
	r.refreshing = true
	r.refreshMutex.Unlock()

    log.Info("Refreshing artist ranks")
	err := r.RefreshRanks()
	if err != nil {
		log.Error("Failed to refresh artist ranks", err)
		// if there's a rate limit from Spotify, avoid repeatedly refreshing
	} else {
		log.Info("Successfully refreshed artist ranks")
	}

	r.lastRefresh = time.Now().Round(0)
	r.refreshMutex.Lock()
	r.refreshing = false
	r.refreshMutex.Unlock()
}

func (r *ArtistRanker) RefreshRanks() error {
	savedTracks, err := r.MusicSvc.GetSavedTracks()
	if err != nil {
		return err
	}
	topTracksLongTerm, err := r.MusicSvc.GetTopTracks(spotify.LongTerm)
	if err != nil {
		return err
	}
	topTracksMediumTerm, err := r.MusicSvc.GetTopTracks(spotify.MediumTerm)
	if err != nil {
		return err
	}
	topTracksShortTerm, err := r.MusicSvc.GetTopTracks(spotify.ShortTerm)
	if err != nil {
		return err
	}
	topArtistsLongTerm, err := r.MusicSvc.GetTopArtists(spotify.LongTerm)
	if err != nil {
		return err
	}
	topArtistsMediumTerm, err := r.MusicSvc.GetTopArtists(spotify.MediumTerm)
	if err != nil {
		return err
	}
	topArtistsShortTerm, err := r.MusicSvc.GetTopArtists(spotify.ShortTerm)
	if err != nil {
		return err
	}
	prevRanks := r.ranks
	r.ranks = make(map[string]rankData, len(topArtistsLongTerm) * 15)

	r.normalizeTrackRanks(topTracksLongTerm)
	r.normalizeTrackRanks(topTracksMediumTerm)
	r.normalizeTrackRanks(topTracksShortTerm)
	r.updateRankForSavedTracks(savedTracks, savedTrackFactor, trackRankCeilingPercent)
	r.updateRankForRankedTracks(topTracksLongTerm, topTrackLongTermFactor, trackRankCeilingPercent)
	r.updateRankForRankedTracks(topTracksMediumTerm, topTrackMediumTermFactor, trackRankCeilingPercent)
	r.updateRankForRankedTracks(topTracksShortTerm, topTrackShortTermFactor, trackRankCeilingPercent)

	r.normalizeArtistRanks(topArtistsLongTerm)
	r.normalizeArtistRanks(topArtistsMediumTerm)
	r.normalizeArtistRanks(topArtistsShortTerm)
	r.updateRankForArtists(topArtistsLongTerm, topArtistLongTermFactor)
	r.updateRankForArtists(topArtistsMediumTerm, topArtistMediumTermFactor)
	r.updateRankForArtists(topArtistsShortTerm, topArtistShortTermFactor)

	err = r.populateRelatedArtistRanks(topArtistsLongTerm)
	if err != nil {
		r.ranks = prevRanks
		return err
	}

	allTracks := slices.Clone(savedTracks)
	for _, t := range topTracksLongTerm {
		allTracks = append(allTracks, t.Track)
	}
	for _, t := range topTracksMediumTerm {
		allTracks = append(allTracks, t.Track)
	}
	for _, t := range topTracksShortTerm {
		allTracks = append(allTracks, t.Track)
	}
	r.populateRelatedForFeaturedArtists(allTracks)

	r.normalizeRanks()

	return nil
}

// from spotify client, tracks are ranked 0 -> N, with top tracks at 0
// convert these to (0, 1], where top tracks have rank 1
func (r ArtistRanker) normalizeTrackRanks(topTracks []spotify.RankedTrack) {
    max := topTracks[len(topTracks)-1].Rank + 1 // add 1 here so even the last place track has rank > 1
	for i := range topTracks {
		track := &topTracks[i]
		track.Rank = (max - track.Rank) / max
	}
}

// from spotify client, artists are ranked 0 -> N, with top artists at 0
// convert these to (0, 1], where top artists have rank 1
func (r ArtistRanker) normalizeArtistRanks(topArtists []spotify.RankedArtist) {
    max := topArtists[len(topArtists)-1].Rank + 1 // add 1 here so even the last place artist has rank > 1
	for i := range topArtists {
		artist := &topArtists[i]
		artist.Rank = (max - artist.Rank) / max
	}
}

func (r *ArtistRanker) normalizeRanks() {
	maxRank := 0.
	for _, rankData := range r.ranks {
		maxRank = max(maxRank, rankData.rank)
	}
	for artist, rankData := range r.ranks {
		rankData.rank /= maxRank
		r.ranks[artist] = rankData
	}
}

func (r *ArtistRanker) updateRankForSavedTracks(tracks []spotify.Track, factor float64, rankCeilPerc float64) {
	tempRanks := make(map[string]float64, len(tracks) * 2)
	maxRank := 0.
	for _, track := range tracks {
		rank := 1.
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
		rankData := r.ranks[key]
		rankData.rank += normalizedRank
		r.ranks[key] = rankData
	}
}

func (r *ArtistRanker) updateRankForRankedTracks(tracks []spotify.RankedTrack, factor float64, rankCeilPerc float64) {
	tempRanks := make(map[string]float64, len(tracks) * 2)
	maxRank := 0.
	for _, track := range tracks {
		rank := track.Rank
		for i, artist := range track.Track.Artists {
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
		rankData := r.ranks[key]
		rankData.rank += normalizedRank
		r.ranks[key] = rankData
	}
}

func (r ArtistRanker) updateRankForArtists(artists []spotify.RankedArtist, weight float64) {
	for _, artist := range artists {
		key :=toKey(artist.Artist.Name)
		rankData := r.ranks[key]
		rankData.rank += artist.Rank * weight
		r.ranks[key] = rankData
	}
}

func (r *ArtistRanker) populateRelatedArtistRanks(artists []spotify.RankedArtist) error {
	log.Infof("Retrieving related artist data for %v artists\n", len(artists))
	// need this so we don't reference increasing ranks of known artists as we iterate
	relatedArtistRanks := make(map[string]float64, len(artists) * 5)

	allArtists := []spotify.Artist{}
	for _, knownArtist := range artists {
		allArtists = append(allArtists, knownArtist.Artist)
	}
	relatedArtistMap, err := r.MusicSvc.GetRelatedArtistsBatch(allArtists)
	if err != nil {
		return err
	}

	for knownArtist, relatedArtists := range relatedArtistMap {
		for _, relatedArtist := range relatedArtists {
			knownArtistData := r.ranks[toKey(knownArtist.Name)]
			calcRankInc := knownArtistData.rank * relatedArtistFactor

			key := toKey(relatedArtist.Name)
			relatedArtistRanks[key] += calcRankInc
			relatedRankData := r.ranks[key]
			if relatedRankData.related == nil {
				relatedRankData.related = []string{}
			}
			relatedRankData.related = append(relatedRankData.related, knownArtist.Name)
			r.ranks[key] = relatedRankData
		}
	}

	for name, rank := range relatedArtistRanks {
		key := toKey(name)
		rankData := r.ranks[key]
		rankData.rank += rank
		r.ranks[key] = rankData
	}
	log.Info("Finished retrieving related artist data")
	return nil
}

func (r *ArtistRanker) populateRelatedForFeaturedArtists(tracks []spotify.Track) {
	related := map[string][]string{}
	for _, track := range tracks {
		primaryArtist := track.Artists[0].Name
		for i, artist := range track.Artists {
			if i == 0 {
				// if it's a primary artist, we will already know enough about them
				related[toKey(primaryArtist)] = nil
			}
			if i > 0 {
				featuredArtist := artist.Name
				key := toKey(featuredArtist)
				if _, exists := related[key]; !exists {
					related[key] = []string{}
				}

				relatedArtists := related[key]
				if relatedArtists != nil {
					relatedArtists = append(relatedArtists, primaryArtist)
					related[key] = relatedArtists
				}
			}
		}
	}
	for artistName, relatedArtists := range related {
		key := toKey(artistName)
		rankData := r.ranks[key]
		if rankData.related == nil {
			rankData.related = []string{}
		}
		rankData.related = relatedArtists
		r.ranks[key] = rankData
	}
}

func toKey(name string) string {
    return strings.ToLower(name)
}

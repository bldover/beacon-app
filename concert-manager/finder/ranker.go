package finder

import (
	"concert-manager/client/lastfm"
	"concert-manager/client/spotify"
	"concert-manager/data"
	"concert-manager/log"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type spotifyService interface {
	GetSavedTracks() ([]spotify.Track, error)
	GetTopTracks(spotify.TimeRange) ([]spotify.Track, error)
	GetTopArtists(spotify.TimeRange) ([]spotify.Artist, error)
}

type lastFmService interface {
	GetSimilarArtists(string) ([]lastfm.Artist, error)
}

type EventRanker struct {
	MusicSvc     spotifyService
	AnalyticsSvc lastFmService
	ranks        map[string]data.ArtistRank
	lastRefresh  time.Time
	refreshing   bool
	refreshMutex sync.Mutex
}

func (r *EventRanker) Rank(event data.EventDetails) data.RankInfo {
	rankInfo := data.RankInfo{Rank: 0, ArtistRanks: map[string]data.ArtistRank{}}

	mainAct := event.Event.MainAct
	if mainAct != nil {
		artistRank := r.rankArtist(*mainAct)
		rankInfo.ArtistRanks[mainAct.Name] = artistRank
		rankInfo.Rank += artistRank.Rank
	}

	for _, opener := range event.Event.Openers {
		artistRank := r.rankArtist(opener)
		rankInfo.ArtistRanks[opener.Name] = artistRank
		rankInfo.Rank += artistRank.Rank
	}

	rankInfo.Recommendation = string(ToRecLevel(rankInfo.Rank))
	log.Debugf("Ranked event %v", rankInfo)
	event.Ranks = &rankInfo
	return rankInfo
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
	relatedArtistFactor  = 0.15
)

var rankTTL, _ = time.ParseDuration("168h")

type track struct {
	TmId    string
	Title   string
	Artists []artist
	Rank    float64
}

type artist struct {
	TmId string
	Name string
	Rank float64
}

func (r *EventRanker) rankArtist(artist data.Artist) data.ArtistRank {
	if r.lastRefresh.IsZero() {
		r.DoRefresh()
	} else if time.Since(r.lastRefresh) > rankTTL {
		go r.DoRefresh()
	}

	rank := r.ranks[toKey(artist.Name)]
	log.Debug("Ranked artist", rank)
	return rank
}

func (r *EventRanker) DoRefresh() {
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

func (r *EventRanker) RefreshRanks() error {
	spotifyTracks, err := r.MusicSvc.GetSavedTracks()
	if err != nil {
		return err
	}
	savedTracks := mapSpotifyTracks(spotifyTracks)
	spotifyTracks, err = r.MusicSvc.GetTopTracks(spotify.LongTerm)
	if err != nil {
		return err
	}
	topTracksLongTerm := mapSpotifyTracks(spotifyTracks)
	spotifyTracks, err = r.MusicSvc.GetTopTracks(spotify.MediumTerm)
	if err != nil {
		return err
	}
	topTracksMediumTerm := mapSpotifyTracks(spotifyTracks)
	spotifyTracks, err = r.MusicSvc.GetTopTracks(spotify.ShortTerm)
	if err != nil {
		return err
	}
	topTracksShortTerm := mapSpotifyTracks(spotifyTracks)
	spotifyArtists, err := r.MusicSvc.GetTopArtists(spotify.LongTerm)
	if err != nil {
		return err
	}
	topArtistsLongTerm := mapSpotifyArtists(spotifyArtists)
	spotifyArtists, err = r.MusicSvc.GetTopArtists(spotify.MediumTerm)
	if err != nil {
		return err
	}
	topArtistsMediumTerm := mapSpotifyArtists(spotifyArtists)
	spotifyArtists, err = r.MusicSvc.GetTopArtists(spotify.ShortTerm)
	if err != nil {
		return err
	}
	topArtistsShortTerm := mapSpotifyArtists(spotifyArtists)

	prevRanks := r.ranks
	r.ranks = make(map[string]data.ArtistRank, len(topArtistsLongTerm)*15)

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

	r.normalizeRanks()
	r.logRanks()

	return nil
}

func mapSpotifyTracks(rawTracks []spotify.Track) []track {
	tracks := []track{}
	for _, rawTrack := range rawTracks {
		tracks = append(tracks, mapSpotifyTrack(rawTrack))
	}
	return tracks
}

func mapSpotifyTrack(rawTrack spotify.Track) track {
	return track{
		TmId:    rawTrack.Id,
		Title:   rawTrack.Title,
		Artists: mapSpotifyArtists(rawTrack.Artists),
	}
}

func mapSpotifyArtists(rawArtists []spotify.Artist) []artist {
	artists := []artist{}
	for _, rawArtist := range rawArtists {
		artists = append(artists, mapSpotifyArtist(rawArtist))
	}
	return artists
}

func mapSpotifyArtist(rawArtist spotify.Artist) artist {
	return artist{
		TmId: rawArtist.Id,
		Name: rawArtist.Name,
	}
}

func mapLastFmArtists(rawArtists []lastfm.Artist) []artist {
	artists := []artist{}
	for _, rawArtist := range rawArtists {
		artists = append(artists, mapLastFmArtist(rawArtist))
	}
	return artists
}

func mapLastFmArtist(rawArtist lastfm.Artist) artist {
	rank, err := strconv.ParseFloat(rawArtist.Rank, 64)
	if err != nil {
		log.Errorf("invalid match value for lastfm artist: %s", rawArtist)
		rank = 0
	}
	return artist{
		Name: rawArtist.Name,
		Rank: float64(rank),
	}
}

// from spotify client, tracks ordered from highest to lowest rank
// convert these to (0, 1], where top tracks have rank 1
func (r *EventRanker) normalizeTrackRanks(topTracks []track) {
	max := float64(len(topTracks))
	for i := range topTracks {
		track := &topTracks[i]
		track.Rank = (max - float64(i)) / max
	}
}

// from spotify client, artists ordered from highest to lowest rank
// convert these to (0, 1], where top artists have rank 1
func (r *EventRanker) normalizeArtistRanks(topArtists []artist) {
	max := float64(len(topArtists))
	for i := range topArtists {
		artist := &topArtists[i]
		artist.Rank = (max - float64(i)) / max
	}
}

func (r *EventRanker) normalizeRanks() {
	maxRank := 0.
	for _, rankData := range r.ranks {
		maxRank = max(maxRank, rankData.Rank)
	}
	for artist, rankData := range r.ranks {
		rankData.Rank /= maxRank
		r.ranks[artist] = rankData
	}
}

func (r *EventRanker) updateRankForSavedTracks(tracks []track, factor float64, rankCeilPerc float64) {
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
		rankData := r.ranks[key]
		rankData.Rank += ceiledRank
		r.ranks[key] = rankData
	}
}

func (r *EventRanker) updateRankForRankedTracks(tracks []track, factor float64, rankCeilPerc float64) {
	tempRanks := make(map[string]float64, len(tracks)*2)
	maxRank := 0.
	for _, track := range tracks {
		rank := track.Rank
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
		rankData.Rank += normalizedRank
		r.ranks[key] = rankData
	}
}

func (r *EventRanker) updateRankForArtists(artists []artist, weight float64) {
	for _, artist := range artists {
		key := toKey(artist.Name)
		rankData := r.ranks[key]
		rankData.Rank += artist.Rank * weight
		r.ranks[key] = rankData
	}
}

func (r *EventRanker) populateRelatedArtistRanks(artists []artist) error {
	log.Infof("Retrieving related artist data for %v artists\n", len(artists))
	// need this so we don't reference increasing ranks of known artists as we iterate
	relatedArtistRanks := make(map[string]float64, len(artists)*5)

	for _, knownArtist := range artists {
		relatedArtistsResp, err := r.AnalyticsSvc.GetSimilarArtists(knownArtist.Name)
		if err != nil {
			log.Errorf("Failed to find related artists for %v, %v", knownArtist, err)
			return err
		}
		relatedArtists := mapLastFmArtists(relatedArtistsResp)
		for _, relatedArtist := range relatedArtists {
			knownArtistData := r.ranks[toKey(knownArtist.Name)]
			calcRankInc := knownArtistData.Rank * relatedArtist.Rank * relatedArtistFactor

			key := toKey(relatedArtist.Name)
			relatedArtistRanks[key] += calcRankInc
			relatedRankData := r.ranks[key]
			if relatedRankData.Related == nil {
				relatedRankData.Related = []string{}
			}
			relatedRankData.Related = append(relatedRankData.Related, knownArtist.Name)
			r.ranks[key] = relatedRankData
		}
	}

	for name, rank := range relatedArtistRanks {
		key := toKey(name)
		rankData := r.ranks[key]
		rankData.Rank += rank
		r.ranks[key] = rankData
	}
	log.Info("Finished retrieving related artist data")
	return nil
}

func toKey(name string) string {
	return strings.ToLower(name)
}

func (r *EventRanker) logRanks() {
	if !log.IsDebug() {
		return
	}

	log.Debug("All artist ranks")
	artists := make([]string, len(r.ranks))
	for artist := range r.ranks {
		artists = append(artists, artist)
	}

	slices.SortFunc(artists, func(a string, b string) int {
		aRank := r.ranks[a].Rank
		bRank := r.ranks[b].Rank
		if aRank < bRank {
			return -1
		} else if aRank > bRank {

			return 1
		} else {
			return 0
		}
	})

	for _, artist := range artists {
		log.Debugf("%s: %v", artist, r.ranks[artist])
	}
}

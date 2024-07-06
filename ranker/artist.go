package ranker

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/spotify"
	"concert-manager/util"
	"time"
)

type ArtistRanker struct {
	MusicSvc musicService
	rankedArtists []string
	ranks map[string]float64
	related map[string][]string // which known artists were used to generate a rank for the input unknown artist
	lastRefresh time.Time
}

type musicService interface {
    GetSavedTracks() ([]spotify.Track, error)
	GetTopTracks(spotify.TimeRange) ([]spotify.RankedTrack, error)
	GetTopArtists(spotify.TimeRange) ([]spotify.RankedArtist, error)
	GetRelatedArtists(spotify.Artist) ([]spotify.Artist, error)
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
	res := data.ArtistRank{Artist: artist, Related: []string{}}
	if artist.Name == "" {
		return res
	}

	if time.Since(r.lastRefresh) > rankTTL {
		log.Info("Refreshing artist ranks")
		err := r.RefreshRanks()
		if err != nil {
			log.Error("Failed to refresh artist ranks", err)
		} else {
			log.Info("Successfully refreshed artist ranks")
			r.lastRefresh = time.Now()
		}
	}

	if rank, ok := r.ranks[artist.Name]; ok {
		res.Rank = rank
		return res
	}

	match := util.SearchStrings(artist.Name, r.rankedArtists, 1, util.ExactTolerance)
	if len(match) > 0 {
		name := match[0]
		log.Debugf("Ranker matched event artist %s with Spotify artist %s\n", artist.Name, name)
		res.Rank = r.ranks[name]
		res.Related = r.related[name]
		return res
	}
	log.Debugf("Ranked artist %v\n", res)
 	return res
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
	r.ranks = make(map[string]float64, len(topArtistsLongTerm) * 15)

	r.normalizeTrackRanks(topTracksLongTerm)
	r.normalizeTrackRanks(topTracksMediumTerm)
	r.normalizeTrackRanks(topTracksShortTerm)
	r.updateRankForSavedTracks(savedTracks, savedTrackFactor, trackRankCeilingPercent)
	r.updateRankForRankedTracks(topTracksLongTerm, topTrackLongTermFactor, trackRankCeilingPercent)
	r.updateRankForRankedTracks(topTracksMediumTerm, topTrackMediumTermFactor, trackRankCeilingPercent)
	r.updateRankForRankedTracks(topTracksShortTerm, topTrackShortTermFactor, trackRankCeilingPercent)

	r.normalizeArtistRanker(topArtistsLongTerm)
	r.normalizeArtistRanker(topArtistsMediumTerm)
	r.normalizeArtistRanker(topArtistsShortTerm)
	r.updateRankForArtists(topArtistsLongTerm, topArtistLongTermFactor)
	r.updateRankForArtists(topArtistsMediumTerm, topArtistMediumTermFactor)
	r.updateRankForArtists(topArtistsShortTerm, topArtistShortTermFactor)

	err = r.populateRelatedArtistRanks(topArtistsLongTerm)
	if err != nil {
		r.ranks = prevRanks
		return err
	}

	r.normalizeRanks()

	r.rankedArtists = make([]string, 0, len(r.ranks))
	for k := range r.ranks {
		r.rankedArtists = append(r.rankedArtists, k)
	}
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
func (r ArtistRanker) normalizeArtistRanker(topArtists []spotify.RankedArtist) {
    max := topArtists[len(topArtists)-1].Rank + 1 // add 1 here so even the last place artist has rank > 1
	for i := range topArtists {
		artist := &topArtists[i]
		artist.Rank = (max - artist.Rank) / max
	}
}

func (r *ArtistRanker) normalizeRanks() {
	maxRank := 0.
	for _, rank := range r.ranks {
		maxRank = max(maxRank, rank)
	}
	for artist, rank := range r.ranks {
		r.ranks[artist] = rank / maxRank
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
			tempRanks[artist.Name] += rank
			maxRank = max(tempRanks[artist.Name], maxRank)
		}
	}

	rankCeiling := maxRank * rankCeilPerc
	for artist, rank := range tempRanks {
		adjRank := min(rank, rankCeiling)
		normalizedRank := (adjRank / rankCeiling) * factor
		r.ranks[artist] += normalizedRank
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
			tempRanks[artist.Name] += rank
			maxRank = max(tempRanks[artist.Name], maxRank)
		}
	}

	rankCeiling := maxRank * rankCeilPerc
	for artist, rank := range tempRanks {
		adjRank := min(rank, rankCeiling)
		normalizedRank := (adjRank / rankCeiling) * factor
		r.ranks[artist] += normalizedRank
	}
}

func (r ArtistRanker) updateRankForArtists(artists []spotify.RankedArtist, weight float64) {
	for _, artist := range artists {
		r.ranks[artist.Artist.Name] += artist.Rank * weight
	}
}

func (r *ArtistRanker) populateRelatedArtistRanks(artists []spotify.RankedArtist) error {
	log.Infof("Retrieving related artist data for %v artists\n", len(artists))
	// need this so we don't reference increasing ranks of known artists as we iterate
	relatedArtistRanks := make(map[string]float64, len(artists) * 5)

	r.related = make(map[string][]string, len(artists) * 5)
	for _, artist := range artists {
		relatedArtists, err := r.MusicSvc.GetRelatedArtists(artist.Artist)
		if err != nil {
			return err
		}
		for _, relatedArtist := range relatedArtists {
			rank := r.ranks[artist.Artist.Name] * relatedArtistFactor
			relatedArtistRanks[relatedArtist.Name] += rank
			if _, ok := r.related[relatedArtist.Name]; !ok {
				r.related[relatedArtist.Name] = []string{}
			}
			r.related[relatedArtist.Name] = append(r.related[relatedArtist.Name], artist.Artist.Name)
		}
	}
	for name, rank := range relatedArtistRanks {
		r.ranks[name] = r.ranks[name] + rank
	}
	log.Info("Finished retrieving related artist data")
	return nil
}

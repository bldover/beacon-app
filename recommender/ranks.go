package recommender

import (
	"concert-manager/spotify"
	"strings"
)

type Recommender struct {
	spotify spotifyClient
}

type spotifyClient interface {
    GetSavedTracks() ([]spotify.Track, error)
	GetTopTracks(spotify.TimeRange) ([]spotify.RankedTrack, error)
	GetTopArtists(spotify.TimeRange) ([]spotify.RankedArtist, error)
	GetRelatedArtists(spotify.Artist) ([]spotify.Artist, error)
}

func NewRecommender() *Recommender {
	client := spotify.NewClient()
	return &Recommender{spotify: client}
}

type RecommendedArtist struct {
    Name string
	Related string
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

func (r *Recommender) RankArtists() (map[RecommendedArtist]float64, error) {
	savedTracks, err := r.spotify.GetSavedTracks()
	if err != nil {
		return nil, err
	}
	topTracksLongTerm, err := r.spotify.GetTopTracks(spotify.LongTerm)
	if err != nil {
		return nil, err
	}
	topTracksMediumTerm, err := r.spotify.GetTopTracks(spotify.MediumTerm)
	if err != nil {
		return nil, err
	}
	topTracksShortTerm, err := r.spotify.GetTopTracks(spotify.ShortTerm)
	if err != nil {
		return nil, err
	}
	topArtistsLongTerm, err := r.spotify.GetTopArtists(spotify.LongTerm)
	if err != nil {
		return nil, err
	}
	topArtistsMediumTerm, err := r.spotify.GetTopArtists(spotify.MediumTerm)
	if err != nil {
		return nil, err
	}
	topArtistsShortTerm, err := r.spotify.GetTopArtists(spotify.ShortTerm)
	if err != nil {
		return nil, err
	}
	ranks := make(map[RecommendedArtist]float64, len(topArtistsLongTerm) * 15)

	normalizeTrackRanks(topTracksLongTerm)
	normalizeTrackRanks(topTracksMediumTerm)
	normalizeTrackRanks(topTracksShortTerm)
	updateRankForTracks(ranks, savedTracks, savedTrackFactor, trackRankCeilingPercent)
	updateRankForTracks(ranks, topTracksLongTerm, topTrackLongTermFactor, trackRankCeilingPercent)
	updateRankForTracks(ranks, topTracksMediumTerm, topTrackMediumTermFactor, trackRankCeilingPercent)
	updateRankForTracks(ranks, topTracksShortTerm, topTrackShortTermFactor, trackRankCeilingPercent)

	normalizeArtistRanks(topArtistsLongTerm)
	normalizeArtistRanks(topArtistsMediumTerm)
	normalizeArtistRanks(topArtistsShortTerm)
	updateRankForArtists(ranks, topArtistsLongTerm, topArtistLongTermFactor)
	updateRankForArtists(ranks, topArtistsMediumTerm, topArtistMediumTermFactor)
	updateRankForArtists(ranks, topArtistsShortTerm, topArtistShortTermFactor)

	err = r.populatedRelatedArtistRanks(ranks, topArtistsLongTerm)
	if err != nil {
		return nil, err
	}

	return ranks, nil
}

// from spotify client, tracks are ranked 0 -> N, with top tracks at 0
// convert these to (0, 1], where top tracks have rank 1
func normalizeTrackRanks(topTracks []spotify.RankedTrack) {
    max := topTracks[len(topTracks)-1].Rank + 1 // add 1 here so even the last place track has rank > 1
	for i := range topTracks {
		track := &topTracks[i]
		track.Rank = (max - track.Rank) / max
	}
}

// from spotify client, artists are ranked 0 -> N, with top artists at 0
// convert these to (0, 1], where top artists have rank 1
func normalizeArtistRanks(topArtists []spotify.RankedArtist) {
    max := topArtists[len(topArtists)-1].Rank + 1 // add 1 here so even the last place artist has rank > 1
	for i := range topArtists {
		artist := &topArtists[i]
		artist.Rank = (max - artist.Rank) / max
	}
}

// build weighted artist frequency map then update artist ranks
func updateRankForTracks[T any](ranks map[RecommendedArtist]float64, tracks []T, factor float64, rankCeilPerc float64) {
	tempRanks := make(map[RecommendedArtist]float64, len(tracks) * 2)
	maxRank := 0.
	for _, track := range tracks {
		rank := 1.
		var artists []spotify.Artist
		if t, ok := any(track).(spotify.RankedTrack); ok {
			rank = t.Rank
			artists = t.Track.Artists
		} else {
			artists = any(track).(spotify.Track).Artists
		}

		for i, artist := range artists {
			if i > 0 {
				rank *= featuredArtistFactor
			}
			recArtist := RecommendedArtist{Name: artist.Name}
			tempRanks[recArtist] += rank
			maxRank = max(tempRanks[recArtist], maxRank)
		}
	}

	rankCeiling := maxRank * rankCeilPerc
	for artist, rank := range tempRanks {
		adjRank := min(rank, rankCeiling)
		normalizedRank := (adjRank / rankCeiling) * factor
		ranks[artist] += normalizedRank
	}
}

func updateRankForArtists(ranks map[RecommendedArtist]float64, artists []spotify.RankedArtist, weight float64) {
	for _, artist := range artists {
		ranks[RecommendedArtist{Name: artist.Artist.Name}] += artist.Rank * weight
	}
}

func (r Recommender) populatedRelatedArtistRanks(ranks map[RecommendedArtist]float64, artists []spotify.RankedArtist) error {
	relatedArtistRanks := make(map[string]float64, len(artists) * 5)
	relatedArtistRefs := make(map[string][]string, len(artists) * 5)
	for _, artist := range artists {
		relatedArtists, err := r.spotify.GetRelatedArtists(artist.Artist)
		if err != nil {
			return err
		}
		for _, relatedArtist := range relatedArtists {
			rank := ranks[RecommendedArtist{Name: artist.Artist.Name}] * relatedArtistFactor
			relatedArtistRanks[relatedArtist.Name] += rank
			if _, ok := relatedArtistRefs[relatedArtist.Name]; !ok {
				relatedArtistRefs[relatedArtist.Name] = []string{}
			}
			relatedArtistRefs[relatedArtist.Name] = append(relatedArtistRefs[relatedArtist.Name], artist.Artist.Name)
		}
	}
	for name, rank := range relatedArtistRanks {
		relatedArtists := formatRelatedArtists(relatedArtistRefs[name])
		existingRank := ranks[RecommendedArtist{Name: name}]
		ranks[RecommendedArtist{Name: name, Related: relatedArtists}] = existingRank + rank
		delete(ranks, RecommendedArtist{Name: name})
	}
	return nil
}

func formatRelatedArtists(refs []string) string {
	sb := strings.Builder{}
	for i, ref := range refs {
		sb.WriteString(ref)
		if i != len(refs) - 1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

package ranker

import (
	"concert-manager/domain"
	"concert-manager/external"
	"concert-manager/log"
	"slices"
)

type RankCalculator struct {
	MusicSvc       spotifyService
	ArtistProvider artistProvider
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

func (calc *RankCalculator) CalculateRanks() (map[string]domain.ArtistRank, error) {
	rawSavedTracks, err := calc.MusicSvc.SavedTracks()
	if err != nil {
		return nil, err
	}
	rawTopTracksLongTerm, err := calc.MusicSvc.TopTracks(external.LongTerm)
	if err != nil {
		return nil, err
	}
	rawTopTracksMediumTerm, err := calc.MusicSvc.TopTracks(external.MediumTerm)
	if err != nil {
		return nil, err
	}
	rawTopTracksShortTerm, err := calc.MusicSvc.TopTracks(external.ShortTerm)
	if err != nil {
		return nil, err
	}
	rawTopArtistsLongTerm, err := calc.MusicSvc.TopArtists(external.LongTerm)
	if err != nil {
		return nil, err
	}
	rawTopArtistsMediumTerm, err := calc.MusicSvc.TopArtists(external.MediumTerm)
	if err != nil {
		return nil, err
	}
	rawTopArtistsShortTerm, err := calc.MusicSvc.TopArtists(external.ShortTerm)
	if err != nil {
		return nil, err
	}

	savedTracks := mapTracks(rawSavedTracks)
	topTracksLongTerm := mapTracks(rawTopTracksLongTerm)
	topTracksMediumTerm := mapTracks(rawTopTracksMediumTerm)
	topTracksShortTerm := mapTracks(rawTopTracksShortTerm)

	topArtistsLongTerm := mapArtists(rawTopArtistsLongTerm)
	topArtistsMediumTerm := mapArtists(rawTopArtistsMediumTerm)
	topArtistsShortTerm := mapArtists(rawTopArtistsShortTerm)

	ranks := make(map[string]domain.ArtistRank, len(topArtistsLongTerm)*15)

	calc.normalizeTrackRanks(topTracksLongTerm)
	calc.normalizeTrackRanks(topTracksMediumTerm)
	calc.normalizeTrackRanks(topTracksShortTerm)
	calc.updateRankForSavedTracks(ranks, savedTracks, savedTrackFactor, trackRankCeilingPercent)
	calc.updateRankForRankedTracks(ranks, topTracksLongTerm, topTrackLongTermFactor, trackRankCeilingPercent)
	calc.updateRankForRankedTracks(ranks, topTracksMediumTerm, topTrackMediumTermFactor, trackRankCeilingPercent)
	calc.updateRankForRankedTracks(ranks, topTracksShortTerm, topTrackShortTermFactor, trackRankCeilingPercent)

	calc.normalizeArtistRanks(topArtistsLongTerm)
	calc.normalizeArtistRanks(topArtistsMediumTerm)
	calc.normalizeArtistRanks(topArtistsShortTerm)
	calc.updateRankForArtists(ranks, topArtistsLongTerm, topArtistLongTermFactor)
	calc.updateRankForArtists(ranks, topArtistsMediumTerm, topArtistMediumTermFactor)
	calc.updateRankForArtists(ranks, topArtistsShortTerm, topArtistShortTermFactor)

	err = calc.populateSimilarArtistRanks(ranks, topArtistsLongTerm)
	if err != nil {
		return nil, err
	}

	calc.normalizeRanks(ranks)
	calc.logRanks(ranks)

	return ranks, nil
}

func (calc *RankCalculator) normalizeTrackRanks(topTracks []rankedTrack) {
	max := float64(len(topTracks))
	for i := range topTracks {
		track := &topTracks[i]
		track.rank = (max - float64(i)) / max
	}
}

func (calc *RankCalculator) normalizeArtistRanks(topArtists []rankedArtist) {
	max := float64(len(topArtists))
	for i := range topArtists {
		artist := &topArtists[i]
		artist.rank = (max - float64(i)) / max
	}
}

func (calc *RankCalculator) normalizeRanks(ranks map[string]domain.ArtistRank) {
	maxRank := 0.
	for _, rankData := range ranks {
		maxRank = max(maxRank, rankData.Rank)
	}
	for artist, rankData := range ranks {
		rankData.Rank /= maxRank
		ranks[artist] = rankData
	}
}

func (calc *RankCalculator) updateRankForSavedTracks(ranks map[string]domain.ArtistRank, tracks []rankedTrack, factor float64, rankCeilPerc float64) {
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
		rankData := ranks[key]
		rankData.Rank += ceiledRank
		ranks[key] = rankData
	}
}

func (calc *RankCalculator) updateRankForRankedTracks(ranks map[string]domain.ArtistRank, tracks []rankedTrack, factor float64, rankCeilPerc float64) {
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
		rankData := ranks[key]
		rankData.Rank += normalizedRank
		ranks[key] = rankData
	}
}

func (calc *RankCalculator) updateRankForArtists(ranks map[string]domain.ArtistRank, artists []rankedArtist, weight float64) {
	for _, artist := range artists {
		key := toKey(artist.Name)
		rankData := ranks[key]
		rankData.Rank += artist.rank * weight
		ranks[key] = rankData
	}
}

func (calc *RankCalculator) populateSimilarArtistRanks(ranks map[string]domain.ArtistRank, artists []rankedArtist) error {
	log.Infof("Retrieving similar artist data for %v artists\n", len(artists))
	similarArtistRanks := make(map[string]float64, len(artists)*5)

	for _, knownArtist := range artists {
		similarArtists, err := calc.ArtistProvider.SimilarArtists(knownArtist.Name)
		if err != nil {
			log.Errorf("Failed to find similar artists for %v, %v", knownArtist, err)
			return err
		}
		for _, similarArtist := range similarArtists {
			if similarArtist.Name == "" {
				continue
			}
			knownArtistData := ranks[toKey(knownArtist.Name)]
			calcRankInc := knownArtistData.Rank * similarArtist.Rank * similarArtistFactor

			key := toKey(similarArtist.Name)
			similarArtistRanks[key] += calcRankInc
			similarRankData := ranks[key]
			if similarRankData.Related == nil {
				similarRankData.Related = []string{}
			}
			similarRankData.Related = append(similarRankData.Related, knownArtist.Name)
			ranks[key] = similarRankData
		}
	}

	for name, rank := range similarArtistRanks {
		key := toKey(name)
		rankData := ranks[key]
		rankData.Rank += rank
		ranks[key] = rankData
	}
	log.Info("Finished retrieving similar artist data")
	return nil
}

func (calc *RankCalculator) logRanks(ranks map[string]domain.ArtistRank) {
	if !log.IsDebug() {
		return
	}

	log.Debug("Logging all artist ranks")
	artists := make([]string, len(ranks))
	for artist := range ranks {
		artists = append(artists, artist)
	}

	slices.SortFunc(artists, func(a string, b string) int {
		aRank := ranks[a].Rank
		bRank := ranks[b].Rank
		if aRank < bRank {
			return -1
		} else if aRank > bRank {
			return 1
		} else {
			return 0
		}
	})

	for _, artist := range artists {
		log.Debugf("%s: %v", artist, ranks[artist])
	}
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
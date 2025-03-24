package ranker

import (
	"concert-manager/data"
	"concert-manager/log"
)

type EventRanker struct {
	ArtistRanker *ArtistRanker
}

func (r *EventRanker) Rank(event data.EventDetails) data.EventRank {
	eventRank := data.EventRank{Event: event, ArtistRanks: []data.ArtistRank{}}
	artistRanking := r.ArtistRanker.Rank(event.Event.MainAct)
	eventRank.Rank = artistRanking.Rank
	eventRank.ArtistRanks = append(eventRank.ArtistRanks, artistRanking)

	for _, opener := range event.Event.Openers {
		artistRanking := r.ArtistRanker.Rank(opener)
		eventRank.Rank += artistRanking.Rank
		eventRank.ArtistRanks = append(eventRank.ArtistRanks, artistRanking)
	}

	if eventRank.Rank == 0 {
		eventRank.Rank += RankGenre(event.EventGenre)
	}
	log.Debugf("Ranked event %v", eventRank)
	return eventRank
}

var genreRanks = map[string]float64{
	"Adult Alternative Pop/Rock":   0.5,
	"Adult Contemporary":           0,
	"African":                      0,
	"Afro-Beat":                    0,
	"Alternative":                  0.5,
	"Alternative Country":          0,
	"Alternative Rap":              0.5,
	"Alternative Rock":             0.5,
	"Americana":                    0,
	"Ballads/Romantic":             0,
	"Bluegrass":                    0,
	"Blues":                        0,
	"British Pop":                  0,
	"Chamber Music":                0,
	"Christian Rap":                0,
	"Christian Rock":               0,
	"Classical":                    0,
	"Club Dance":                   0,
	"Colombia":                     0,
	"Comedy":                       0,
	"Contemporary Christian":       0,
	"Country":                      0,
	"Country Soul":                 0,
	"Dance-Rock":                   0,
	"Dance/Electronic":             0,
	"Dancehall Reggae":             0,
	"Death Metal/ Black Metal":     1,
	"Death Metal/Black Metal":      1,
	"Disco":                        0,
	"Dream Pop":                    0.5,
	"Electronic":                   0,
	"Electronic Pop":               0.5,
	"Emo":                          1,
	"Experimental Rock":            1,
	"Folk":                         0,
	"Folk Rock":                    0.5,
	"Garage Rock":                  1,
	"Gospel":                       0,
	"Goth Metal":                   1,
	"Hard Rock":                    1,
	"Hardcore Punk":                1,
	"Heavy Metal":                  1,
	"Hip-Hop/Rap":                  0.5,
	"Hobby/Special Interest Expos": 0,
	"India & Pakistan":             0,
	"Indie Folk":                   0,
	"Indie Pop":                    0.5,
	"Indie Rock":                   1,
	"Japanese Pop":                 0,
	"Jazz":                         0,
	"K-Pop":                        0,
	"Latin":                        1,
	"Latin Electronica":            0,
	"Latin Pop":                    1,
	"Latin Rap":                    1,
	"Medieval/Renaissance":         0,
	"Metal":                        1,
	"Mexican Grupero":              1,
	"Middle Age":                   0,
	"Music":                        0,
	"Oldies & Classics":            0,
	"Other":                        0,
	"Pop":                          0.5,
	"Pop Punk":                     0.5,
	"Pop Rock":                     0.5,
	"Pop-Soul":                     0,
	"Post-Hardcore":                1,
	"Post-Punk":                    1,
	"Psychedelic":                  0.5,
	"Punk":                         0.5,
	"R&B":                          0,
	"Reggae":                       0,
	"Religious":                    0,
	"Rock":                         1,
	"Rock & Roll":                  1,
	"Singer-Songwriter":            0,
	"Soul":                         0,
	"Symphonic":                    0,
	"Synth Pop":                    0.5,
	"Undefined":                    0,
	"Urban":                        0.5,
	"World":                        0.5,
}

const genreWeight = 0.1

func RankGenre(genre string) float64 {
	if rank, exists := genreRanks[genre]; exists {
		return rank * genreWeight
	}
	return 0
}

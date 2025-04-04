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

	log.Debugf("Ranked event %v", eventRank)
	return eventRank
}

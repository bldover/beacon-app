package ranker

import (
	"concert-manager/domain"
	"concert-manager/log"
)

type EventRanker struct {
	Cache *ArtistRankCache
}

func (r *EventRanker) Rank(event domain.EventDetails) domain.RankInfo {
	rankInfo := domain.RankInfo{ArtistRanks: map[string]domain.ArtistRank{}}

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
	log.Debugf("Ranked event %v, %v", event, rankInfo)
	return rankInfo
}

func (r *EventRanker) rankArtist(artist domain.Artist) domain.ArtistRank {
	rank := r.Cache.Rank(artist)
	log.Debugf("Ranked artist %s, %v", artist.Name, rank)
	return rank
}

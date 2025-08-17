package lastfm

import (
	"concert-manager/external"
	"concert-manager/log"
	"strconv"
)

type similarArtistResponse struct {
	Similar similarArtistList `json:"similarartists"`
}

type similarArtistList struct {
	Artists []artist `json:"artist"`
}

type artist struct {
	Name string `json:"name"`
	Rank string `json:"match"`
}

func (c *Client) SimilarArtists(artist string) ([]external.RankedArtist, error) {
	queryParams := map[string]any{}
	queryParams["method"] = "artist.getsimilar"
	queryParams["artist"] = artist

	request := requestEntity{queryParams}
	response := &similarArtistResponse{}
	err := c.call(request, response)
	if err != nil {
		return nil, err
	}

	related := []external.RankedArtist{}
	if response.Similar.Artists == nil || len(response.Similar.Artists) == 0 {
		log.Infof("No related artists found for %s", artist)
		return related, nil
	}

	for _, artist := range response.Similar.Artists {
		related = append(related, mapArtist(artist))
	}

	return related, nil
}

func mapArtist(rawArtist artist) external.RankedArtist {
	rank, err := strconv.ParseFloat(rawArtist.Rank, 64)
	if err != nil {
		log.Errorf("invalid match value for lastfm artist: %s", rawArtist)
		rank = 0
	}
	return external.RankedArtist{
		Name: rawArtist.Name,
		Rank: float64(rank),
	}
}

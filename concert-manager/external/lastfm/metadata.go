package lastfm

import (
	"concert-manager/external"
	"concert-manager/log"
	"strings"
)

type infoResponse struct {
	Error   int                 `json:"error"`
	Message string              `json:"message"`
	Artist  *artistInfoResponse `json:"artist"`
}

type artistInfoResponse struct {
	Name    string            `json:"name"`
	MBID    string            `json:"mbid"`
	Url     string            `json:"url"`
	Similar similarArtistList `json:"similar"`
	Tags    tagsResponse      `json:"tags"`
}

type tagsResponse struct {
	Tag []tag `json:"tag"`
}

type tag struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (c *Client) ArtistInfoById(mbid string) (external.ArtistInfo, error) {
	log.Info("Request to get LastFM artist details for ID", mbid)
	queryParams := map[string]any{}
	queryParams["method"] = "artist.getinfo"
	queryParams["mbid"] = mbid

	request := requestEntity{queryParams}
	response := &infoResponse{}
	err := c.call(request, response)
	if err != nil {
		return external.ArtistInfo{}, err
	}
	if response.Error == 6 {
		return external.ArtistInfo{}, external.NotFoundError{Message: "LastFM artist not found for MBID " + mbid}
	}

	artistInfo := mapArtistInfo(*response.Artist)
	log.Info("Retrieved LastFM details", artistInfo)
	return artistInfo, nil
}

func (c *Client) SearchByName(name string) (external.ArtistInfo, error) {
	log.Info("Request to get LastFM artist details for name", name)
	queryParams := map[string]any{}
	queryParams["method"] = "artist.getinfo"
	queryParams["artist"] = name

	request := requestEntity{queryParams}
	response := &infoResponse{}
	info := external.ArtistInfo{}

	err := c.call(request, response)
	if err != nil {
		return external.ArtistInfo{}, err
	}
	if response.Error == 6 {
		return external.ArtistInfo{}, external.NotFoundError{Message: "LastFM artist not found for name " + name}
	}

	info = mapArtistInfo(*response.Artist)
	log.Info("Retrieved LastFm artist details", info)
	return info, nil
}

func mapArtistInfo(r artistInfoResponse) external.ArtistInfo {
	genres := []string{}
	for _, tag := range r.Tags.Tag {
		// lots of odd things in the artist tags since they are user added, but "seen live" is very common
		if !strings.EqualFold(tag.Name, "seen live") {
			genres = append(genres, tag.Name)
		}
	}

	return external.ArtistInfo{
		Name:   r.Name,
		Genres: genres,
		Id:     r.MBID,
	}
}

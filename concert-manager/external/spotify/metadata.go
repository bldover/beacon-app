package spotify

import (
	"concert-manager/external"
	"concert-manager/log"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const artistsPath = "/artists"

func (c *Client) ArtistInfoById(artistId string) (external.ArtistInfo, error) {
	log.Info("Request to get Spotify artist details for ID", artistId)

	artistsUrl := baseUrl + artistsPath + "/" + artistId
	request := RequestEntity{artistsUrl, nil}
	response := &artist{}

	err := c.call(http.MethodGet, request, response)
	if err != nil {
		if errorResponse, ok := err.(errorResponse); ok {
			switch errorResponse.ErrorDetails.Status {
			case http.StatusNotFound:
				fallthrough
			case http.StatusBadRequest:
				return external.ArtistInfo{}, external.NotFoundError{Message: "Spotify artist not found for ID " + artistId}
			default:
				return external.ArtistInfo{}, errors.New("failed to retrieve Spotify artist details: " + errorResponse.Error())
			}
		}
		return external.ArtistInfo{}, err
	}

	artistInfo := external.ArtistInfo{
		Name:   response.Name,
		Genres: response.Genres,
		Id:     response.Id,
	}

	log.Info("Retrieved Spotify details", artistInfo)
	return artistInfo, nil
}

const searchPath = "/search"

func (c *Client) SearchByName(name string) (external.ArtistInfo, error) {
	log.Info("Request to get Spotify artist details for name", name)
	url := baseUrl + searchPath
	queryParams := map[string]any{
		"q":     name,
		"type":  "artist",
		"limit": 2, // for some reason, Spotify sometimes returns the wrong artist with limit 1
	}

	request := RequestEntity{url, queryParams}
	response := &artistSearchResponse{}
	info := external.ArtistInfo{}

	err := c.call(http.MethodGet, request, response)
	if err != nil {
		if errorResponse, ok := err.(errorResponse); ok {
			switch errorResponse.ErrorDetails.Status {
			case http.StatusNotFound:
				fallthrough
			case http.StatusBadRequest:
				return external.ArtistInfo{}, external.NotFoundError{Message: "Spotify artist not found for name " + name}
			default:
				return external.ArtistInfo{}, fmt.Errorf("failed to retrieve Spotify artist details for name %s: %v", name, errorResponse.Error())
			}
		}
		return external.ArtistInfo{}, err
	}

	artists := response.Artists.Items
	if len(artists) == 0 {
		return info, fmt.Errorf("empty spotify search results for %v", name)
	}
	var artist artist
	for _, artistResp := range response.Artists.Items {
		if strings.EqualFold(artistResp.Name, name) {
			artist = artistResp
		}
	}

	if artist.Name == "" {
		return info, fmt.Errorf("no artist found in spotify search results for %v", name)
	}

	info = external.ArtistInfo{Name: name, Id: artists[0].Id, Genres: artists[0].Genres}
	log.Info("Retrieved Spotify artist details", info)
	return info, nil
}

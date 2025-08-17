package ticketmaster

import (
	"concert-manager/domain"
	"concert-manager/log"
	"concert-manager/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type tmResponse struct {
	Links struct {
		Next struct {
			URL string `json:"href"`
		} `json:"next"`
	} `json:"_links"`
	Data struct {
		Events []tmEventResponse `json:"events"`
	} `json:"_embedded"`
	PageInfo struct {
		EventCount int `json:"totalElements"`
		Page       int `json:"number"`
	} `json:"page"`
}

type tmEventResponse struct {
	EventName string `json:"name"`
	Id        string `json:"id"`
	Dates     struct {
		Start struct {
			Date string `json:"localDate"`
		} `json:"start"`
		Status struct {
			Code string `json:"code"`
		} `json:"status"`
	} `json:"dates"`
	Classification []tmGenreResponse `json:"classification"`
	Details        struct {
		Venues []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			City struct {
				Name string `json:"Name"`
			} `json:"city"`
			State struct {
				Name string `json:"name"`
			} `json:"state"`
		} `json:"venues"`
		Artists []struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Links struct {
				Wiki []struct {
					URL string `json:"url"`
				} `json:"wiki"`
				Spotify []struct {
					URL string `json:"url"`
				} `json:"spotify"`
				MusicBrainz []struct {
					ID string `json:"id"`
				} `json:"musicbrainz"`
			} `json:"externalLinks"`
			Classification []tmGenreResponse `json:"classifications"`
		} `json:"attractions"`
	} `json:"_embedded"`
}

type tmGenreResponse struct {
	Genre struct {
		Name string `json:"name"`
	} `json:"genre"`
	Subgenre struct {
		Name string `json:"name"`
	} `json:"subGenre"`
}

type errorResponse struct {
	Fault struct {
		Details struct {
			Code string `json:"errorcode"`
		} `json:"detail"`
	} `json:"fault"`
}

func toResponse(body io.Reader) (*tmResponse, error) {
	var resp tmResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		errMsg := fmt.Sprintf("failed to parse ticketmaster response: %v", err)
		return nil, errors.New(errMsg)
	}
	return &resp, nil
}

func toErrorResponse(body io.Reader) (*errorResponse, error) {
	var resp errorResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		errMsg := fmt.Sprintf("failed to parse ticketmaster error response: %v", err)
		return nil, errors.New(errMsg)
	}
	return &resp, nil
}

func toSpotifyId(url string) string {
	if url == "" {
		return ""
	}

	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

func parseEventDetails(event *tmEventResponse) (*domain.EventDetails, error) {
	eventName := event.EventName
	artistDetails := event.Details.Artists
	if eventName == "" && len(artistDetails) == 0 {
		return nil, errors.New("no event name or artists")
	}

	mainAct := domain.Artist{}
	if len(artistDetails) != 0 {
		mainActDetails := artistDetails[0]
		mainAct.Name = mainActDetails.Name
		if len(mainActDetails.Classification) != 0 {
			mainAct.Genres.Ticketmaster = []string{getGenre(mainActDetails.Classification[0])}
		}
		mainAct.ID.Ticketmaster = mainActDetails.ID
		if len(mainActDetails.Links.Spotify) > 0 {
			mainAct.ID.Spotify = toSpotifyId(mainActDetails.Links.Spotify[0].URL)
		}
		if len(mainActDetails.Links.MusicBrainz) > 0 {
			mainAct.ID.MusicBrainz = mainActDetails.Links.MusicBrainz[0].ID
		}
	}

	openers := []domain.Artist{}
	if len(artistDetails) > 1 {
		for _, openerDetails := range artistDetails[1:] {
			if openerDetails.Name == "" {
				log.Infof("no opener artist name in event %+v", event)
				continue
			}
			opener := domain.Artist{
				Name: openerDetails.Name,
			}
			if len(openerDetails.Classification) != 0 {
				opener.Genres.Ticketmaster = []string{getGenre(openerDetails.Classification[0])}
			}
			opener.ID.Ticketmaster = openerDetails.ID
			if len(openerDetails.Links.Spotify) > 0 {
				opener.ID.Spotify = toSpotifyId(openerDetails.Links.Spotify[0].URL)
			}
			if len(openerDetails.Links.MusicBrainz) > 0 {
				opener.ID.MusicBrainz = openerDetails.Links.MusicBrainz[0].ID
			}

			openers = append(openers, opener)
		}
	}

	eventGenre := ""
	if len(event.Classification) != 0 {
		eventGenre = getGenre(event.Classification[0])
	}

	venue := domain.Venue{}
	if len(event.Details.Venues) != 0 {
		venueDetails := event.Details.Venues[0]
		venue.Name = venueDetails.Name
		venue.City = venueDetails.City.Name
		venue.State = venueDetails.State.Name
	}

	dateRaw := event.Dates.Start.Date
	date, err := time.Parse(dateFmt, dateRaw)
	if err != nil {
		errMsg := fmt.Sprintf("unable to parse event date %s", dateRaw)
		return nil, errors.New(errMsg)
	}

	eventDetails := domain.EventDetails{
		Name:       eventName,
		EventGenre: eventGenre,
		Event: domain.Event{
			MainAct: &mainAct,
			Openers: openers,
			Venue:   venue,
			Date:    util.Date(date),
			ID:      domain.ID{Ticketmaster: event.Id},
		},
	}

	if event.Dates.Status.Code == "cancelled" || eventDetails.Event.MainAct.Name == "Test artist" {
		return &eventDetails, eventCancelledError{"Event has been cancelled"}
	}
	return &eventDetails, nil
}

func getGenre(genres tmGenreResponse) string {
	subGenre := genres.Subgenre.Name
	genre := genres.Genre.Name
	switch {
	case subGenre != "" && subGenre != "Undefined" && subGenre != "Other":
		return subGenre
	case genre != "" && genre != "Undefined" && genre != "Other":
		return genre
	default:
		return ""
	}
}

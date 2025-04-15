package ticketmaster

import (
	"concert-manager/data"
	"concert-manager/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
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
	Prices []struct {
		MinPrice float64 `json:"min"`
	} `json:"priceRanges"`
	Ticketing struct {
		InclusivePricing struct {
			Enabled bool `json:"enabled"`
		} `json:"allInclusivePricing"`
	} `json:"ticketing"`
	Classification []tmGenreResponse `json:"classification"`
	Details struct {
		Venues []struct {
			Name string `json:"name"`
			City struct {
				Name string `json:"Name"`
			} `json:"city"`
			State struct {
				Name string `json:"name"`
			} `json:"state"`
		} `json:"venues"`
		Artists []struct {
			Name  string `json:"name"`
			Links struct {
				Wiki []struct {
					URL string `json:"url"`
				} `json:"wiki"`
				Spotify []struct {
					URL string `json:"url"`
				} `json:"spotify"`
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

func parseEventDetails(event *tmEventResponse) (*data.EventDetails, error) {
	eventName := event.EventName
	artistDetails := event.Details.Artists
	if eventName == "" && len(artistDetails) == 0 {
		return nil, errors.New("no event name or artists")
	}

	mainAct := data.Artist{}
	if len(artistDetails) != 0 {
		mainActDetails := artistDetails[0]
		mainAct.Name = mainActDetails.Name
		if len(mainActDetails.Classification) != 0 {
			mainAct.Genres.TmGenres = []string{getGenre(mainActDetails.Classification[0])}
		}
	}

	openers := []data.Artist{}
	if len(artistDetails) > 1 {
		for _, openerDetails := range artistDetails[1:] {
			if openerDetails.Name == "" {
				return nil, errors.New("no opener artist name")
			}
			opener := data.Artist{
				Name: openerDetails.Name,
			}
			if len(openerDetails.Classification) != 0 {
				opener.Genres.TmGenres = []string{getGenre(openerDetails.Classification[0])}
			}
			openers = append(openers, opener)
		}
	}

	eventGenre := ""
	if len(event.Classification) != 0 {
		eventGenre = getGenre(event.Classification[0])
	}

	price := ""
	if len(event.Prices) == 0 {
		price = "Unknown"
	} else {
 		price = strconv.FormatFloat(event.Prices[0].MinPrice, 'f', 2, 64)
		if !event.Ticketing.InclusivePricing.Enabled {
			price += " + fees"
		}
	}

	venue := data.Venue{}
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

	eventDetails := data.EventDetails{
		Name:  eventName,
		Price: price,
		EventGenre: eventGenre,
		Event: data.Event{
			MainAct: &mainAct,
			Openers: openers,
			Venue:   venue,
			Date:    util.Date(date),
			TmId:    event.Id,
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

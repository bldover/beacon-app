package finder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	Dates     struct {
		Start struct {
			Date string `json:"localDate"`
		} `json:"start"`
	} `json:"dates"`
	Prices []struct {
		MinPrice float64 `json:"min"`
	} `json:"priceRanges"`
	Ticketing struct {
		InclusivePricing struct {
			Enabled bool `json:"enabled"`
		} `json:"allInclusivePricing"`
	} `json:"ticketing"`
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
			Classification []struct {
				Genre struct {
					Name string `json:"name"`
				} `json:"genre"`
				Subgenre struct {
					Name string `json:"name"`
				} `json:"subGenre"`
			} `json:"classifications"`
		} `json:"attractions"`
	} `json:"_embedded"`
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

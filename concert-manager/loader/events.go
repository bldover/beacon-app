package loader

import (
	"bufio"

	"concert-manager/domain"
	"concert-manager/log"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

const minColumns = 7

type eventCache interface {
	AddSavedEvent(domain.Event) (*domain.Event, error)
}

type EventLoader struct {
	Cache eventCache
}

// requires a UTF-8 encoded CSV file
func (l *EventLoader) Upload(ctx context.Context, file io.ReadCloser) (int, error) {
	log.Debug("Starting processing event file upload")
	scanner := bufio.NewScanner(file)
	first := true
	events := []domain.Event{}
	for scanner.Scan() {
		if first {
			first = false
			continue
		}
		line := scanner.Text()
		log.Debugf("Parsed event line: %s", line)
		event, err := toEvent(line)
		if err != nil {
			log.Errorf("Error while parsing event: %v", err)
			return 0, fmt.Errorf("unable to convert line to event: %s, %v", line, err)
		}
		log.Debugf("Converted input to event %v", event)
		events = append(events, event)
	}

	hasErr := false
	successCount := 0
	for i, event := range events {
		log.Debugf("Starting upload for event %v", event)
		if _, err := l.Cache.AddSavedEvent(event); err != nil {
			log.Errorf("Failed to add event at row %d, %+v, %v", i+2, event, err)
			hasErr = true
		} else {
			successCount++
			log.Debugf("Event successfully uploaded %v", event)
		}
	}

	log.Infof("Successfully uploaded %d event rows", successCount)
	log.Errorf("Failed to upload %d event rows", len(events)-successCount)
	if hasErr {
		return successCount, errors.New("failed to add at least one row. check logs for more details")
	}
	return successCount, nil
}

func toEvent(row string) (domain.Event, error) {
	parts := strings.Split(row, ",")
	if len(parts) < minColumns {
		return domain.Event{}, errors.New("not enough columns in row")
	}

	mainAct := domain.Artist{
		Name: strings.TrimSpace(parts[0]),
	}
	mainAct.Genres.Ticketmaster = []string{}
	genres := strings.Split(strings.TrimSpace(parts[1]), ";")
	mainAct.Genres.User = append(mainAct.Genres.User, genres...)
	if !mainAct.Populated() {
		return domain.Event{}, errors.New("invalid main act")
	}

	date := strings.TrimSpace(parts[2])
	venue := domain.Venue{
		Name:  strings.TrimSpace(parts[3]),
		City:  strings.TrimSpace(parts[4]),
		State: strings.TrimSpace(parts[5]),
	}
	if !venue.Populated() {
		return domain.Event{}, errors.New("invalid venue")
	}
	purchased := strings.TrimSpace(parts[6]) == "TRUE"

	openers := []domain.Artist{}
	i, j := 7, 8
	for i < len(parts) && j < len(parts) {
		opener := domain.Artist{
			Name: strings.TrimSpace(parts[i]),
		}
		mainAct.Genres.Ticketmaster = []string{}
		genres := strings.Split(strings.TrimSpace(parts[j]), ";")
		mainAct.Genres.User = append(mainAct.Genres.User, genres...)
		if !opener.Populated() {
			return domain.Event{}, errors.New("invalid opener")
		}
		if !opener.Populated() {
			break
		}
		openers = append(openers, opener)
		i, j = i+2, j+2
	}

	event := domain.Event{
		MainAct:   &mainAct,
		Openers:   openers,
		Venue:     venue,
		Date:      date,
		Purchased: purchased,
	}

	return event, nil
}

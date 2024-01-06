package svc

import (
	"bufio"
	"bytes"
	"concert-manager/data"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

const minColumns = 6

type EventCreator interface {
	AddEventRecursive(context.Context, Event) error
}

type Loader struct {
	EventCreator EventCreator
}

func (l *Loader) Upload(ctx context.Context, file io.ReadCloser) (int, error) {
	scanner := bufio.NewScanner(file)
	first := true
	events := []data.Event{}
 	for scanner.Scan() {
		if first {
			first = false
			continue
		}
		line := scanner.Text()
		event, err := toEvent(convertSpecialChars(line))
		if err != nil {
			return 0, fmt.Errorf("unable to convert line to event: %s, %v", line, err)
		}
		events = append(events, event)
	}

	hasErr := false
	successCount := 0
	for i, event := range events {
		fmt.Printf("parsed event %+v to upload", event)
		if err := l.EventCreator.AddEventRecursive(ctx, event); err != nil {
			log.Printf("Failed to add event at row %d, %+v, %v", i+2, event, err)
			hasErr = true
		} else {
			successCount++
			fmt.Println("Event successfully uploaded")
		}
	}

	if hasErr {
		return successCount, errors.New("failed to add at least one row. check logs for more details")
	}
	return successCount, nil
}

func toEvent(row string) (Event, error) {
	parts := strings.Split(row, ",")
	if len(parts) < minColumns {
		return data.Event{}, errors.New("not enough columns in row")
	}

	mainAct := data.Artist{
		Name: strings.TrimSpace(parts[0]),
		Genre: strings.TrimSpace(parts[1]),
	}
	if mainAct.Invalid() {
		return Event{}, errors.New("invalid main act")
	}
	date := strings.TrimSpace(parts[2])
	venue := data.Venue{
		Name: strings.TrimSpace(parts[3]),
		City: strings.TrimSpace(parts[4]),
		State: strings.TrimSpace(parts[5]),
	}
	if !venue.Populated() {
		return Event{}, errors.New("invalid venue")
	}
	openers := []data.Artist{}
	i, j := 6, 7
	for i < len(parts) && j < len(parts) {
		opener := Artist{
			Name: strings.TrimSpace(parts[i]),
			Genre: strings.TrimSpace(parts[j]),
		}
		if opener.Invalid() {
			return Event{}, errors.New("invalid opener")
		}
		if !opener.Populated() {
			break
		}
		openers = append(openers, opener)
		i, j = i + 2, j + 2
	}

	event := Event{
		MainAct: mainAct,
		Openers: openers,
		Venue: venue,
		Date: date,
	}

	return event, nil
}

func convertSpecialChars(input string) string {
    var buf bytes.Buffer
    for _, r := range input {
        if unicode.IsControl(r) {
            fmt.Fprintf(&buf, "\\u%04X", r)
        } else {
            fmt.Fprintf(&buf, "%c", r)
        }
    }
    return buf.String()
}

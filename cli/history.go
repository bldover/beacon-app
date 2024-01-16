package cli

import (
	"bufio"
	"concert-manager/data"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const pageSize = 10

type History struct {
	Events   *[]data.Event
	ParentScreen Screen
	AddHistScreen Screen
	DeleteHistScreen HistoryDeleter
	page     int
}

type HistoryDeleter interface {
	Screen
	AddContext(int, int)
}

const (
	nextPage = iota + 1
	prevPage
	gotoPage
	addEvent
	deleteEvent
)

func (h History) numPages() int {
    return int(math.Ceil(float64(len(*h.Events)) / float64(pageSize)))
}

func (h History) Title() string {
	return "Concert History"
}

func (h History) Data() string {
	if len(*h.Events) == 0 {
		return "No concerts found"
	}

	var data strings.Builder
	pageIndicator := fmt.Sprintf("Page %d/%d\n", h.page+1, h.numPages())
	data.WriteString(pageIndicator)
	startEvent := (h.page * pageSize)
	endEvent := startEvent + pageSize
	if endEvent > len(*h.Events) {
		endEvent = len(*h.Events)
	}
	for i := startEvent; i < endEvent; i++ {
		data.WriteString(formatEvent((*h.Events)[i]))
	}
	return data.String()
}

func formatEvent(e data.Event) string {
	format := "%v @ %s\n\tArtists: %s\n\tGenres: %s\n"
	location := fmt.Sprintf("%s, %s, %s", e.Venue.Name, e.Venue.City, e.Venue.State)
	artists := []string{}
	genres := []string{}
	if e.MainAct.Populated() {
		artists = append(artists, e.MainAct.Name)
		genres = append(genres, e.MainAct.Genre)
	}
	for _, artist := range e.Openers {
		if artist.Populated() {
			if !slices.Contains(artists, artist.Name) {
				artists = append(artists, artist.Name)
			}
			if !slices.Contains(genres, artist.Genre) {
				genres = append(genres, artist.Genre)
			}
		}
	}
	artistStr := strings.Join(artists, ", ")
	genreStr := strings.Join(genres, ", ")
	return fmt.Sprintf(format, e.Date, location, artistStr, genreStr)
}

func (h History) Actions() []string {
	return []string{
		"Next Page",
		"Prev Page",
		"Goto Page",
		"Add Event",
		"Delete Event",
	}
}

func (h *History) NextScreen(i int) Screen {
	switch i {
    case nextPage:
		if (h.page + 1) < h.numPages() {
			h.page++
		}
		return h
	case prevPage:
		if h.page > 0 {
			h.page--
		}
	case gotoPage:
		h.handleGoto()
	case addEvent:
		return h.AddHistScreen
	case deleteEvent:
		h.DeleteHistScreen.AddContext(pageSize * h.page, pageSize)
		return h.DeleteHistScreen
	}
	return h
}

func (h *History) handleGoto() {
	fmt.Println("Enter page number:")
	reader := bufio.NewReader(os.Stdin)
	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading selection, try again:")
			continue
		}

		page, err := strconv.Atoi(in[:len(in) - 1])
		if err != nil {
			fmt.Println("Invalid option, try again:")
			continue
		}
		if page > h.numPages() || page < 1 {
			fmt.Println("Invalid option, try again:")
			continue
		}
		h.page = page - 1
		break
	}
}

func (h History) Parent() Screen {
    return h.ParentScreen
}

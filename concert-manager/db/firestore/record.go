package firestore

import (
	"concert-manager/domain"
	"concert-manager/log"
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const recordCollection = "records"

var recordFields = []string{"Name", "Artist", "Year", "Signed", "ID"}

type (
	RecordClient struct {
		Connection *Firestore
	}

	RecordEntity struct {
		Name   string
		Artist ArtistEntity
		Year   int
		Signed bool
		ID     string
	}
)

func (c *RecordClient) Add(ctx context.Context, record domain.Record) (string, error) {
	log.Debug("Attempting to add record", record)
	existingRecord, err := c.findDocRef(ctx, record.ID)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error occurred while checking if record %v already exists, %v", record, err)
		return "", err
	}
	if existingRecord.Exists() {
		log.Debugf("Skipping adding record because it already exists %+v, %v", record, existingRecord.Ref.ID)
		return existingRecord.Ref.ID, nil
	}

	artistEntity := toArtistEntity(record.Artist)
	recordEntity := RecordEntity{
		Name:   record.Name,
		Artist: artistEntity,
		Year:   record.Year,
		Signed: record.Signed,
	}

	records := c.Connection.Client.Collection(recordCollection)
	docRef, _, err := records.Add(ctx, recordEntity)
	if err != nil {
		log.Errorf("Failed to add new record %+v, %v", record, err)
		return "", err
	}
	log.Infof("Created new record %+v", docRef.ID)
	return docRef.ID, nil
}

func (c *RecordClient) Update(ctx context.Context, record domain.Record) error {
	log.Debug("Attempting to update record", record)
	recordDoc, err := c.Connection.Client.Collection(recordCollection).Doc(record.ID).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Errorf("Error while updating record %+v, %v", record, err)
		return err
	}
	if !recordDoc.Exists() {
		log.Errorf("Record does not exist in update for %+v, %v", record, err)
		return err
	}

	artistEntity := toArtistEntity(record.Artist)
	recordEntity := RecordEntity{
		Name:   record.Name,
		Artist: artistEntity,
		Year:   record.Year,
		Signed: record.Signed,
		ID:     record.ID,
	}

	_, err = recordDoc.Ref.Set(ctx, recordEntity)
	if err != nil {
		log.Errorf("Failed to update record %+v, %v", record, err)
		return err
	}
	log.Info("Successfully updated record", record)
	return nil
}

func (c *RecordClient) Delete(ctx context.Context, id string) error {
	log.Debug("Attempting to delete record", id)
	recordDoc, err := c.Connection.Client.Collection(recordCollection).Doc(id).Get(ctx)
	if err != nil {
		log.Error("Error while deleting record", id, err)
		return err
	}
	_, err = recordDoc.Ref.Delete(ctx)
	if err != nil {
		log.Error("Failed to delete record", id, err)
		return err
	}
	log.Info("Successfully deleted record", id)
	return nil
}

func (c *RecordClient) FindAll(ctx context.Context) ([]domain.Record, error) {
	log.Debug("Finding all records")
	recordDocs, err := c.Connection.Client.Collection(recordCollection).
		Select(recordFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		log.Error("Error while finding all records,", err)
		return nil, err
	}

	records := []domain.Record{}
	for _, r := range recordDocs {
		records = append(records, toRecord(r))
	}
	log.Debugf("Found %d records", len(records))
	return records, nil
}

func toRecord(doc *firestore.DocumentSnapshot) domain.Record {
	data := doc.Data()
	record := domain.Record{
		ID: doc.Ref.ID,
	}

	if name, ok := data["Name"].(string); ok {
		record.Name = name
	}
	if year, ok := data["Year"].(int64); ok {
		record.Year = int(year)
	}
	if signed, ok := data["Signed"].(bool); ok {
		record.Signed = signed
	}
	if artistData, ok := data["Artist"].(map[string]any); ok {
		record.Artist = toArtistFromMap(artistData, "")
	}

	return record
}

func toArtistEntity(artist domain.Artist) ArtistEntity {
	genreEntity := GenreInfoEntity{
		Spotify:      artist.Genres.Spotify,
		LastFm:       artist.Genres.LastFm,
		Ticketmaster: artist.Genres.Ticketmaster,
		User:         artist.Genres.User,
	}
	idEntity := ArtistIDEntity{
		Primary:      artist.ID.Primary,
		Spotify:      artist.ID.Spotify,
		Ticketmaster: artist.ID.Ticketmaster,
		MusicBrainz:  artist.ID.MusicBrainz,
	}
	return ArtistEntity{Name: artist.Name, Genres: genreEntity, ID: idEntity}
}

func toArtistFromMap(data map[string]any, docID string) domain.Artist {
	artist := domain.Artist{
		Genres: domain.GenreInfo{
			Spotify:      []string{},
			LastFm:       []string{},
			Ticketmaster: []string{},
			User:         []string{},
		},
	}

	if name, ok := data["Name"].(string); ok {
		artist.Name = name
	}

	if ids, ok := data["ID"].(map[string]any); ok {
		if primary, ok := ids["Primary"].(string); ok {
			artist.ID.Primary = primary
		}
		if spotify, ok := ids["Spotify"].(string); ok {
			artist.ID.Spotify = spotify
		}
		if tm, ok := ids["Ticketmaster"].(string); ok {
			artist.ID.Ticketmaster = tm
		}
		if mb, ok := ids["MusicBrainz"].(string); ok {
			artist.ID.MusicBrainz = mb
		}
	}

	if genres, ok := data["Genres"].(map[string]any); ok {
		if spotify, ok := genres["Spotify"].([]any); ok {
			for _, g := range spotify {
				if gs, ok := g.(string); ok {
					artist.Genres.Spotify = append(artist.Genres.Spotify, gs)
				}
			}
		}
		if lastFm, ok := genres["LastFm"].([]any); ok {
			for _, g := range lastFm {
				if gs, ok := g.(string); ok {
					artist.Genres.LastFm = append(artist.Genres.LastFm, gs)
				}
			}
		}
		if user, ok := genres["User"].([]any); ok {
			for _, g := range user {
				if gs, ok := g.(string); ok {
					artist.Genres.User = append(artist.Genres.User, gs)
				}
			}
		}
	}

	return artist
}

func (c *RecordClient) findDocRef(ctx context.Context, id string) (*firestore.DocumentSnapshot, error) {
	if id == "" {
		return &firestore.DocumentSnapshot{}, nil
	}
	return c.Connection.Client.Collection(recordCollection).Doc(id).Get(ctx)
}

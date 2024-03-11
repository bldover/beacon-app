package firestore

import (
	"concert-manager/log"
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
)

type Firestore struct {
	Client *firestore.Client
}

func Setup() (*Firestore, error) {
	projectID := os.Getenv("PROJ_ID")
	if projectID == "" {
		return nil, errors.New("PROJ_ID environment variable must be set")
	}
	log.Debug("Connecting to Firestore in project", projectID)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	fs := Firestore{client}
	return &fs, nil
}

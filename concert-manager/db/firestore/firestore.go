package firestore

import (
	"concert-manager/log"
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
)

const (
	projectIdEnv = "CM_PROJ_ID"
)

type Firestore struct {
	Client *firestore.Client
}

func Setup() (*Firestore, error) {
	projectID := os.Getenv(projectIdEnv)
	if projectID == "" {
		return nil, errors.New("CM_PROJ_ID environment variable must be set")
	}
	log.Debug("Connecting to Firestore in project", projectID)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	fs := Firestore{client}
	log.Info("Successfully initialized database")
	return &fs, nil
}

package db

import (
	"concert-manager/out"
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

type Firestore struct {
	Client *firestore.Client
}

func Setup() (*Firestore, error) {
	projectID := os.Getenv("PROJ_ID")
	out.Debugf("Connecting to Firestore in project %s", projectID)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	fs := Firestore{client}
	return &fs, nil
}

package gcs

import (
	"concert-manager/log"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
)

const albumImageBucketEnv = "CM_ALBUM_IMAGE_BUCKET"

type GCS struct {
	client *storage.Client
	bucket string
}

func Setup() (*GCS, error) {
	bucket := os.Getenv(albumImageBucketEnv)
	if bucket == "" {
		return nil, errors.New("CM_ALBUM_IMAGE_BUCKET environment variable must be set")
	}
	log.Debug("Connecting to GCS bucket", bucket)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	log.Info("Successfully initialized GCS client")
	return &GCS{client: client, bucket: bucket}, nil
}

func (g *GCS) UploadImage(ctx context.Context, content io.Reader, contentType string) (string, error) {
	objectName := generateObjectName(contentType)
	wc := g.client.Bucket(g.bucket).Object(objectName).NewWriter(ctx)
	wc.ContentType = contentType

	if _, err := io.Copy(wc, content); err != nil {
		wc.Close()
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucket, objectName)
	log.Infof("Successfully uploaded image to GCS: %s", url)
	return url, nil
}

func generateObjectName(contentType string) string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b) + extensionFromContentType(contentType)
}

func extensionFromContentType(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

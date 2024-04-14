package server

import (
	"concert-manager/log"
	"concert-manager/spotify"
	"context"
	"fmt"
	"io"
	"net/http"
)

const Port = ":3001"
const maxFileSizeBytes = 100000

func StartServer(l Loader) {
	http.Handle("/upload", &uploadHandler{l})
	http.Handle("/test", &testHandler{})
	http.Handle("/spotify/callback", &spotify.SpotifyAuthHandler{})
	log.Info("Starting server on port", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}

type Loader interface {
    Upload(context.Context, io.ReadCloser) (int, error)
}

type uploadHandler struct {
	loader Loader
}

func (handler *uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		errMsg := fmt.Sprintf("Error while parsing request file %v", err)
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
	}

	log.Info("Received upload request")

    rows, err := handler.loader.Upload(r.Context(), file)

	if err != nil {
		errMsg := fmt.Sprintf("Error occurred during upload processing: %v", err)
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	successMsg := fmt.Sprintf("Successfully uploaded %d rows", rows)
	log.Info("Finished processing upload request")
	fmt.Fprintln(w, successMsg)
}

type testHandler struct {}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("No active test")
}

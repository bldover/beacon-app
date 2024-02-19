package server

import (
	"concert-manager/out"
	"context"
	"fmt"
	"io"
	"net/http"
)

const port = ":3001"
const maxFileSizeBytes = 100000

func StartServer(l Loader) {
	handler := &uploadHandler{l}
	http.Handle("/upload", handler)
	out.Infoln("Starting server on port ", port)
	out.Fatal(http.ListenAndServe(port, nil))
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
		out.Errorln(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
	}

	out.Infoln("Received upload request.")

    rows, err := handler.loader.Upload(r.Context(), file)

	if err != nil {
		errMsg := fmt.Sprintf("Error occurred during upload processing: %v", err)
		out.Errorln(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	successMsg := fmt.Sprintf("Successfully uploaded %d rows", rows)
	out.Infoln("Finished processing upload request")
	fmt.Fprintln(w, successMsg)
}

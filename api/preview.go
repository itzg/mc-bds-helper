package api

import (
	"log"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func GetLatestPreview(w http.ResponseWriter, _ *http.Request) {
	url, err := lookupLatestVersion(typePreview)
	if err != nil {
		log.Printf("E: %s", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(err.statusCode)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(url))
}

package api

import (
	"github.com/itzg/mc-bds-helper/lookup"
	"log"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func GetLatestPreview(w http.ResponseWriter, _ *http.Request) {
	url, err := lookup.LatestVersion(lookup.TypePreview)
	if err != nil {
		log.Printf("E: %s", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(err.StatusCode)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(url))
}

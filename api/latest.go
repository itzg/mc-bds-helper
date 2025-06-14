package api

import (
	"github.com/itzg/mc-bds-helper/internal"
	"log"
	"net/http"
)

func GetLatest(w http.ResponseWriter, _ *http.Request) {
	url, err := internal.LookupLatestVersion(internal.TypeRelease)
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

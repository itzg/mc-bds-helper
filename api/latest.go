package api

import (
	"fmt"
	"github.com/itzg/mc-bds-helper/lookup"
	"log"
	"net/http"
	"time"
)

const CacheAge = 60 * time.Minute

func GetLatest(w http.ResponseWriter, _ *http.Request) {
	url, err := lookup.LatestVersion(lookup.TypeRelease)
	if err != nil {
		writeError(w, err)
		return
	}

	writeUrlResponse(w, url)
}

func writeError(w http.ResponseWriter, err *lookup.LookupError) {
	log.Printf("E: %s", err)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(err.StatusCode)
	_, _ = w.Write([]byte(err.Error()))
}

func writeUrlResponse(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "text/plain")
	writeCacheHeaders(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(url))
}

func writeCacheHeaders(w http.ResponseWriter) {
	cacheAgeStr := fmt.Sprintf("max-age=%d", CacheAge/time.Second)
	w.Header().Set("Cache-Control", cacheAgeStr)
	w.Header().Set("CDN-Cache-Control", cacheAgeStr)
	w.Header().Set("Vercel-CDN-Cache-Control", cacheAgeStr)
}

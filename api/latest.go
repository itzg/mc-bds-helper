package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultCacheDuration = time.Hour
	downloadLinksUrl     = "https://net-secondary.web.minecraft-services.net/api/v1.0/download/links"
	typeRelease          = "serverBedrockLinux"
)

var (
	cachedArchiveUrl = ""
	cacheUntil       time.Time
)

type DownloadLinksResponse struct {
	Result struct {
		Links []struct {
			DownloadType string `json:"downloadType"`
			DownloadUrl  string `json:"downloadUrl"`
		} `json:"links"`
	} `json:"result"`
}

func GetLatest(w http.ResponseWriter, _ *http.Request) {
	if cachedArchiveUrl == "" || time.Now().After(cacheUntil) {
		var err *lookupError
		cachedArchiveUrl, err = lookupLatestVersion()
		if err != nil {
			log.Printf("E: %s", err)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(err.statusCode)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		cacheDuration := loadCacheDuration()
		cacheUntil = time.Now().Add(cacheDuration)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(cachedArchiveUrl))
}

func loadCacheDuration() time.Duration {
	var cacheDuration time.Duration
	cacheDurationStr := os.Getenv("CACHE_DURATION")
	if cacheDurationStr != "" {
		var parseErr error
		cacheDuration, parseErr = time.ParseDuration(cacheDurationStr)
		if parseErr != nil {
			cacheDurationStr = ""
		}
	}
	if cacheDurationStr == "" {
		cacheDuration = defaultCacheDuration
	}
	return cacheDuration
}

type lookupError struct {
	message    string
	wrapped    error
	statusCode int
}

func newLookupError(message string, wrapped error, statusCode int) *lookupError {
	return &lookupError{message: message, wrapped: wrapped, statusCode: statusCode}
}

func (e *lookupError) Unwrap() error {
	return e.wrapped
}

func (e *lookupError) Error() string {
	return e.message
}

func lookupLatestVersion() (string, *lookupError) {
	resp, err := http.Get(downloadLinksUrl)

	if err != nil {
		return "", newLookupError("http issue", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		return "", newLookupError("failed to lookup latest version", nil, resp.StatusCode)
	}

	DownloadLinksReponse := DownloadLinksResponse{}
	err = json.NewDecoder(resp.Body).Decode(&DownloadLinksReponse)
	if err != nil {
		return "", newLookupError("failed to decode response", err, http.StatusInternalServerError)
	}

	for _, link := range DownloadLinksReponse.Result.Links {
		if link.DownloadType == typeRelease {
			return link.DownloadUrl, nil
		}
	}

	return "", newLookupError("failed to find release link", nil, http.StatusInternalServerError)
}

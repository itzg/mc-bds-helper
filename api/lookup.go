package api

import (
	"encoding/json"
	cache "github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"time"
)

var versionsCache = cache.New(60*time.Minute, 120*time.Minute)

const (
	downloadLinksUrl = "https://net-secondary.web.minecraft-services.net/api/v1.0/download/links"
	typeRelease      = "serverBedrockLinux"
	typePreview      = "serverBedrockLinuxPreview"
)

type DownloadLinksResponse struct {
	Result struct {
		Links []struct {
			DownloadType string `json:"downloadType"`
			DownloadUrl  string `json:"downloadUrl"`
		} `json:"links"`
	} `json:"result"`
}

func lookupLatestVersion(downloadType string) (string, *lookupError) {
	url, present := versionsCache.Get(downloadType)
	if present {
		log.Print("Using cached version", downloadType)
		return url.(string), nil
	}

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
		if link.DownloadType == downloadType {
			versionsCache.Set(link.DownloadType, link.DownloadUrl, cache.DefaultExpiration)
			return link.DownloadUrl, nil
		}
	}

	return "", newLookupError("failed to find release link", nil, http.StatusInternalServerError)
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

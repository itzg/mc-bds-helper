package api

import (
	"net/http/httptest"
	"testing"
)

func TestGetLatestVersion(t *testing.T) {
	resp := &httptest.ResponseRecorder{}
	GetLatest(resp, httptest.NewRequest("GET", "/api/latest-version", nil))

	defer resp.Result().Body.Close()
	if resp.Result().StatusCode != 200 {
		t.Errorf("Response code was %d", resp.Result().StatusCode)
	}
}

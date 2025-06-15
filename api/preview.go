package api

import (
	"github.com/itzg/mc-bds-helper/lookup"
	"net/http"
)

//goland:noinspection GoUnusedExportedFunction
func GetLatestPreview(w http.ResponseWriter, _ *http.Request) {
	url, err := lookup.LatestVersion(lookup.TypePreview)
	if err != nil {
		writeError(w, err)
		return
	}

	writeUrlResponse(w, url)
}

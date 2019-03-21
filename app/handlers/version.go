package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	versionResponse struct {
		Version string `json:"version"`
	}

	versionHandler struct {
		version string
		logger  *log.Logger
	}
)

func (v *versionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "" {
		contentType = "application/json"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	response := versionResponse{
		Version: v.version,
	}

	responseBody, err := json.Marshal(response)

	if err != nil {
		v.logger.Println(err)
		_, err = w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
		if err != nil {
			v.logger.Println(err)
		}
	}

	logRequest(v.logger, r, commonLogFormat, http.StatusOK, string(responseBody))

	json.NewEncoder(w).Encode(response)
	return
}

// VersionHandler allocates and returns a new versionHandler to report version.
func VersionHandler(version string, logger *log.Logger) http.Handler {
	return &versionHandler{
		version: version,
		logger:  logger,
	}
}

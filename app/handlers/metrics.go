package handlers

import (
	"fmt"
	"github.com/ahmedkamals/foo-protocol-proxy/app/analysis"
	"log"
	"net/http"
)

type (
	// metricsHandler acts as an interface for the metrics data that should be exported over HTTP.
	metricsHandler struct {
		analyzer *analysis.Analyzer
		logger   *log.Logger
	}
)

// MetricsHandler allocates and returns a new metricsHandler to report stats.
func MetricsHandler(analyzer *analysis.Analyzer, logger *log.Logger) http.Handler {
	return &metricsHandler{
		analyzer: analyzer,
		logger:   logger,
	}
}

func (m *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "" {
		contentType = "application/json"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	//w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	//w.Header().Set("Content-Disposition", "attachment; filename=\"results.json\"")

	body, err := m.analyzer.Report()

	if err != nil {
		m.logger.Println(err)
		_, err = w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
		if err != nil {
			m.logger.Println(err)
		}
	}

	_, err = w.Write([]byte(body))
	if err != nil {
		m.logger.Println(err)
	}

	logRequest(m.logger, r, combinedLogFormat, http.StatusOK, body)
}

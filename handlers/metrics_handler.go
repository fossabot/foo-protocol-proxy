package handlers

import (
	"fmt"
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"log"
	"net/http"
	"time"
)

type (
	// MetricsHandler acts as an interface for the metrics data that should be exported over HTTP.
	MetricsHandler struct {
		analyzer *analysis.Analyzer
		logger   *log.Logger
	}
)

// CommonLogFormat 125.125.125.125 - dsmith [10/Oct/1999:21:15:05 +0500] "GET /index.html HTTP/1.0" 200 1043
// For more details: https://en.wikipedia.org/wiki/Common_Log_Format
const CommonLogFormat string = "%s %s %s [%s] \"%s %s %s\" %d %d"

// NewMetricsHandler allocates and returns a new MetricsHandler to report stats.
func NewMetricsHandler(analyzer *analysis.Analyzer, logger *log.Logger) http.Handler {
	return &MetricsHandler{
		analyzer: analyzer,
		logger:   logger,
	}
}

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	response, err := m.analyzer.Report()

	if err != nil {
		m.logger.Println(err)
		_, err = w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
		if err != nil {
			m.logger.Println(err)
		}
	}

	_, err = w.Write([]byte(response))
	if err != nil {
		m.logger.Println(err)
	}

	m.logger.Printf(
		CommonLogFormat,
		r.RemoteAddr,
		"-",
		"-",
		time.Now().Format("01/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.Path,
		r.Proto,
		http.StatusOK,
		len(response),
	)
}

package handlers

import (
	"foo-protocol-proxy/analysis"
	"log"
	"net/http"
)

type (
	MetricsHandler struct {
		analyzer *analysis.Analyzer
	}
)

func NewMetricsHandler(analyzer *analysis.Analyzer) http.Handler {
	return &MetricsHandler{
		analyzer: analyzer,
	}
}

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "" {
		contentType = "application/json"
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	//w.Header().Set("Content-Disposition", "attachment; filename=\"results.json\"")

	data, err := m.analyzer.Report()

	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(data))
}

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type (
	// HealthHandler acts as an interface for server health over HTTP.
	HealthHandler struct {
		logger *log.Logger
	}
)

// CombinedLogFormat host rfc931 username date:time request statuscode bytes referrer user_agent cookie
// e.g. 125.125.125.125 - dsmith [10/Oct/1999:21:15:05 +0500] "GET /index.html HTTP/1.0" 200 1043 "http://www.ibm.com/" "Mozilla/4.05 [en] (WinNT; I)" "USERID=CustomerA;IMPID=01234"
// For more details: http://publib.boulder.ibm.com/tividd/td/ITWSA/ITWSA_info45/en_US/HTML/guide/c-logs.html#combined
const CombinedLogFormat string = "%s %s %s [%s] \"%s %s %s\" %d %d %s %s %s"

// NewHealthHandler allocates and returns a new HealthHandler to report health.
func NewHealthHandler(logger *log.Logger) http.Handler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "" {
		contentType = "application/json"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	data := map[string]string{"status": "OK!"}
	response, err := json.Marshal(data)

	if err != nil {
		h.logger.Println(err)
		_, err = w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
		if err != nil {
			h.logger.Println(err)
		}
	}

	_, err = w.Write(response)
	if err != nil {
		h.logger.Println(err)
	}
	//json.NewEncoder(w).Encode(map[string]string{"status": "OK!"})

	h.logRequest(CombinedLogFormat, http.StatusOK, r, response)
}

func (h *HealthHandler) logRequest(format string, statusCode int, r *http.Request, response []byte) {

	referer := "-"
	if r.Referer() != "" {
		referer = fmt.Sprintf("\"%s\"", r.Referer())
	}

	userAgent := "-"
	if r.UserAgent() != "" {
		userAgent = fmt.Sprintf("\"%s\"", r.UserAgent())
	}

	cookiesRaw := "-"

	if len(r.Cookies()) > 0 {

		var cookiesValues []string
		for _, cookie := range r.Cookies() {
			cookiesValues = append(cookiesValues, cookie.String())
		}

		cookiesRaw = fmt.Sprintf("\"%s\"", strings.Join(cookiesValues, ";"))
	}

	h.logger.Printf(
		format,
		r.RemoteAddr,
		"-",
		"-",
		time.Now().Format("01/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.Path,
		r.Proto,
		statusCode,
		len(response),
		referer,
		userAgent,
		cookiesRaw,
	)
}

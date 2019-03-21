package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	baseFormat = "%s %s %s [%s] \"%s %s %s\""
	// combinedLogFormat host rfc931 username date:time request status code bytes referrer user_agent cookie
	// e.g. 125.125.125.125 - akamal [10/Oct/1999:21:15:05 +0500] "GET /index.html HTTP/1.0" 200 1043 "http://www.ibm.com/" "Mozilla/4.05 [en] (WinNT; I)" "USERID=CustomerA;IMPID=01234"
	// For more details: http://publib.boulder.ibm.com/tividd/td/ITWSA/ITWSA_info45/en_US/HTML/guide/c-logs.html#combined
	combinedLogFormat string = baseFormat + " %d %d %s %s %s"
	// commonLogFormat 125.125.125.125 - akamal [10/Oct/1999:21:15:05 +0500] "GET /index.html HTTP/1.0" 200 1043
	// For more details: https://en.wikipedia.org/wiki/Common_Log_Format
	commonLogFormat string = baseFormat + " %d %d"
	// headLogFormat for request that sets only header, and does not return any body.
	headLogFormat = baseFormat + " %s\n"

	dateFormat string = "01/Jan/2006:15:04:05 -0700"
)

// LoggingHandler returns an HTTP handler.
func LoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventLogger := log.New(os.Stdout, "", 0)
		logRequest(eventLogger, r, headLogFormat, 0, "")
		h.ServeHTTP(w, r)
	})
}

func buildValues(r *http.Request, format string, statusCode int, body string) []interface{} {
	values := []interface{}{
		r.RemoteAddr,
		"-",
		"-",
		time.Now().Format(dateFormat),
		r.Method,
		r.URL.Path,
		r.Proto,
	}

	switch format {
	case combinedLogFormat, commonLogFormat:
		values = append(values, statusCode, len(body))
	case headLogFormat:
		values = append(values, r.UserAgent())
	}

	if format != combinedLogFormat {
		return values
	}

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

	return append(values, referer, userAgent, cookiesRaw)
}

func logRequest(logger *log.Logger, r *http.Request, format string, statusCode int, body string) {
	values := buildValues(r, format, statusCode, body)
	logger.Printf(format, values...)
}

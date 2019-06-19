package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/ahmedkamals/foo-protocol-proxy/app/testingutil"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShouldHandleHealthRoutesCorrectly(t *testing.T) {
	t.Parallel()
	expectedOutput, err := json.Marshal(map[string]string{"version": "test-version"})

	if err != nil {
		t.Error(err)
	}

	expected := string(expectedOutput)
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.Ldate)
	testCases := []testingutil.TestCase{
		{
			ID: "Version route",
			Input: map[string]http.Handler{
				"/version": VersionHandler("test-version", logger),
			},
			Expected: expected,
		},
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/route", nil)

	if err != nil {
		t.Error(err)
	}

	for _, testCase := range testCases {
		input := testCase.Input.(map[string]http.Handler)
		expected := testCase.Expected.(string)

		for _, handler := range input {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Unexpectd error code %d", w.Code)
			}

			if w.Header().Get("content-type") != "application/json" {
				t.Errorf("Unexpectd header %s", w.Header().Get("Content-Type"))
			}

			assert.Equal(t, expected, strings.TrimSpace(w.Body.String()))
		}
	}
}

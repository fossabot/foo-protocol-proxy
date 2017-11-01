package app

import (
	"errors"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/handlers"
	"github.com/ahmedkamals/foo-protocol-proxy/testingutil"
	"net/http"
	"reflect"
	"testing"
)

func TestShouldConfigureRoutesCorrectly(t *testing.T) {
	t.Parallel()
	routes := map[string]http.Handler{
		"/health": handlers.NewHealthHandler(),
	}
	server := getServer(config.Configuration{}, routes)
	mux := http.NewServeMux()

	for route, handler := range routes {
		mux.Handle(route, handler)
	}

	testCases := []testingutil.TestCase{
		{
			ID:           "Configuring routes - error should occur",
			Input:        map[string]http.Handler{},
			Expected:     errors.New("missing routes"),
			ExpectsError: true,
		},
		{
			ID:           "Configuring routes - No error",
			Input:        routes,
			Expected:     mux,
			ExpectsError: false,
		},
	}

	for _, testCase := range testCases {
		input := testCase.Input.(map[string]http.Handler)
		actual, err := server.configureRoutesHandler(input)

		if testCase.ExpectsError {
			expected := testCase.Expected.(error)

			if err.Error() != expected.Error() {
				t.Error(testCase.Format(actual))
			}
			continue
		} else {
			expected := testCase.Expected.(*http.ServeMux)

			if !reflect.DeepEqual(expected, actual) {
				t.Error(testCase.Format(actual))
			}
		}
	}
}

func TestShouldGetServerRoutesCorrectly(t *testing.T) {
	t.Parallel()
	routes := map[string]http.Handler{
		"/health": handlers.NewHealthHandler(),
	}
	server := getServer(config.Configuration{}, routes)

	testCases := []testingutil.TestCase{
		{
			ID:       "Server routes",
			Input:    server,
			Expected: routes,
		},
	}

	for _, testCase := range testCases {
		input := testCase.Input.(*HTTPServer)
		expected := testCase.Expected.(map[string]http.Handler)
		actual := input.getRoutes()

		if !reflect.DeepEqual(expected, actual) {
			t.Error(testCase.Format(actual))
		}
	}
}

func getServer(config config.Configuration, routes map[string]http.Handler) *HTTPServer {
	return NewHTTPServer(config, routes, make(chan error, 1))
}

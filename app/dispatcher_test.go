package app

import (
	"bytes"
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/handlers"
	"github.com/ahmedkamals/foo-protocol-proxy/testingutil"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestShouldGetDispatcherRoutesCorrectly(t *testing.T) {
	t.Parallel()
	analyzer := analysis.NewAnalyzer()
	routes := getRoutes(analyzer)
	saver := getMockedSaver()
	dispatcher := NewDispatcher(config.Configuration{}, analyzer, saver)

	testCases := []testingutil.TestCase{
		{
			ID:       "Dispatcher routes",
			Input:    dispatcher,
			Expected: routes,
		},
	}

	for _, testCase := range testCases {
		input := testCase.Input.(*Dispatcher)
		expected := testCase.Expected.(map[string]http.Handler)
		actual := input.getRoutes()

		assert.Len(t, actual, len(expected))

		for key, val := range actual {
			assert.NotEmpty(t, actual[key])
			assert.Implements(t, (*http.Handler)(nil), val)
		}
	}
}

func TestShouldBlockIndefinitely(t *testing.T) {
	t.Parallel()
	testCases := []testingutil.TestCase{
		{
			ID:       "Connection forwarding",
			Input:    make(chan os.Signal, 1),
			Expected: true,
		},
	}

	configuration := config.Configuration{}
	analyzer := analysis.NewAnalyzer()
	saver := getMockedSaver()
	dispatcher := NewDispatcher(configuration, analyzer, saver)

	for _, testCase := range testCases {
		input := testCase.Input.(chan os.Signal)
		expected := testCase.Expected.(bool)
		input <- os.Interrupt
		actual := dispatcher.blockIndefinitely(input, true)

		assert.Equal(t, expected, actual)
	}
}

func getRoutes(analyzer *analysis.Analyzer) map[string]http.Handler {
	var buf bytes.Buffer
	eventLogger := log.New(&buf, "", log.Ldate)

	return map[string]http.Handler{
		"/metrics":   handlers.NewMetricsHandler(analyzer, eventLogger),
		"/stats":     handlers.NewMetricsHandler(analyzer, eventLogger),
		"/health":    handlers.NewHealthHandler(eventLogger),
		"/heartbeat": handlers.NewHealthHandler(eventLogger),
		"/status":    handlers.NewHealthHandler(eventLogger),
	}
}

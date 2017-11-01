package app

import (
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/handlers"
	"github.com/ahmedkamals/foo-protocol-proxy/persistence"
	"github.com/ahmedkamals/foo-protocol-proxy/testingutil"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestShouldGetDispatcherRoutesCorrectly(t *testing.T) {
	t.Parallel()
	analyzer := analysis.NewAnalyzer()
	routes := map[string]http.Handler{
		"/metrics": handlers.NewMetricsHandler(analyzer),
		"/stats":   handlers.NewMetricsHandler(analyzer),
		"/health":  handlers.NewHealthHandler(),
		"/status":  handlers.NewHealthHandler(),
	}
	dispatcher := NewDispatcher(config.Configuration{}, analyzer, persistence.NewSaver(""))

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

		if !reflect.DeepEqual(expected, actual) {
			t.Error(testCase.Format(actual))
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
	saver := persistence.NewSaver("")
	dispatcher := NewDispatcher(configuration, analyzer, saver)

	for _, testCase := range testCases {
		input := testCase.Input.(chan os.Signal)
		expected := testCase.Expected.(bool)
		input <- os.Interrupt
		actual := dispatcher.blockIndefinitely(input, true)

		if expected != actual {
			t.Error(testCase.Format(actual))
		}
	}
}

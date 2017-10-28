package app

import (
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/persistence"
	"github.com/ahmedkamals/foo-protocol-proxy/testingutil"
	"os"
	"testing"
)

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

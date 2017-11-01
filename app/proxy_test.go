package app

import (
	"errors"
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/persistence"
	"github.com/ahmedkamals/foo-protocol-proxy/testingutil"
	"net"
	"reflect"
	"testing"
)

func TestShouldStartProperly(t *testing.T) {
	t.Parallel()
	listener, err := getMockedListener(":8010")

	if err != nil {
		t.Error(err)
	}

	go listener.Accept()

	proxy := getProxy(config.Configuration{
		Forwarding:   ":8010",
		Listening:    ":8011",
		HTTPAddress:  "0.0.0.0:8001",
		RecoveryPath: "data/recovery.json",
	})

	testCases := []testingutil.TestCase{
		{
			ID:       "Connection forwarding",
			Input:    proxy,
			Expected: "*net.TCPConn",
		},
		{
			ID:           "Connection forwarding - should return error",
			Input:        getProxy(config.Configuration{}),
			Expected:     errors.New(""),
			ExpectsError: true,
		},
	}

	for _, testCase := range testCases {
		input := testCase.Input.(*Proxy)
		actual, err := input.forward()

		if !testCase.ExpectsError {
			expected := testCase.Expected.(string)

			if reflect.TypeOf(actual).String() != expected {
				t.Error(testCase.Format(actual))
			}
		} else if err == nil {
			t.Error(testCase.Format("There should be an error"))
		}
	}

	listener.Close()
}

func getProxy(config config.Configuration) *Proxy {
	return NewProxy(config,
		getMockedAnalyzer(),
		getMockedSaver(),
		make(chan error, 10),
	)
}

func getMockedListener(listeningPort string) (net.Listener, error) {
	return net.Listen("tcp", listeningPort)
}

func getMockedAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{}
}

func getMockedSaver() *persistence.Saver {
	return &persistence.Saver{}
}

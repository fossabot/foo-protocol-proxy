package app

import (
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/persistance"
	"github.com/ahmedkamals/foo-protocol-proxy/testing_util"
	"net"
	"reflect"
	"testing"
)

func TestShouldStartProperly(t *testing.T) {
	listener, err := getMockedListener(":8010")

	if err != nil {
		t.Error(err)
	}

	go listener.Accept()

	proxy := NewProxy(config.Configuration{
		Forwarding:   ":8010",
		Listening:    ":8011",
		HttpAddress:  "0.0.0.0:8001",
		RecoveryPath: "data/recovery.json",
	},
		getMockedAnalyzer(),
		getMockedSaver(),
	)

	testCases := []testing_util.TestCase{
		{
			Id:       "Connection forwarding",
			Input:    proxy,
			Expected: "",
		},
	}

	for _, testCase := range testCases {
		input := testCase.Input.(*Proxy)
		//expected := testCase.Expected
		actual, err := input.forward()

		if err != nil {
			t.Error(err)
		}

		if reflect.TypeOf(actual).String() != "*net.TCPConn" {
			t.Error(testCase.Format(actual))
		}
	}

	listener.Close()
}

func getMockedListener(listeningPort string) (net.Listener, error) {
	return net.Listen("tcp", listeningPort)
}

func getMockedAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{}
}

func getMockedSaver() *persistance.Saver {
	return &persistance.Saver{}
}

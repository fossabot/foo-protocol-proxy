package persistance

import (
	"encoding/json"
	testingUtil "github.com/ahmedkamals/foo-protocol-proxy/testing_util"
	"reflect"
	"testing"
)

func TestShouldUnmarshallCorrectly(t *testing.T) {

	data := `{"index":7,"time_stamp":1505951090,"requests_10s":[1,0,0,0,0,0,0,0,0,0],"responses_10s":[1,0,0,0,0,0,0,0,0,0]}`
	expectedRequests := [10]uint64{}
	expectedResponses := [10]uint64{}
	expectedRequests[0] = 1
	expectedResponses[0] = 1

	testCases := []testingUtil.TestCase{
		{
			Id:       "Should unmarshall recovered data correctly.",
			Input:    []byte(data),
			Expected: NewRecovery(7, 1505951090, expectedRequests, expectedResponses),
		},
	}

	for _, testCase := range testCases {
		input := testCase.Input.([]byte)
		expected := testCase.Expected.(*Recovery)
		actual := NewEmptyRecovery()
		err := json.Unmarshal(input, actual)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Error(testCase.Format(actual))
		}
	}
}

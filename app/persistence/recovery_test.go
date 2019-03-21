package persistence

import (
	"encoding/json"
	testingUtil "github.com/ahmedkamals/foo-protocol-proxy/app/testingutil"
	"github.com/stretchr/testify/assert"
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
			ID:       "Should unmarshall recovered data correctly.",
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

		assert.Equal(t, expected, actual)
	}
}

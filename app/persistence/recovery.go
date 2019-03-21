package persistence

import (
	"encoding/json"
)

type (
	// Recovery type manages data that should be recovered later on.
	Recovery struct {
		Index             uint8      `json:"index"`
		TimeStamp         uint64     `json:"time_stamp"`
		RequestsInTenSec  [10]uint64 `json:"requests_10s"`
		ResponsesInTenSec [10]uint64 `json:"responses_10s"`
	}
)

// NewRecovery allocates and returns a new Recovery
// from a given data
// to track the data that should be recovered later.
func NewRecovery(index uint8, timestamp uint64, requestsInTenSec, responsesInTenSec [10]uint64) *Recovery {
	return &Recovery{
		Index:             index,
		TimeStamp:         timestamp,
		RequestsInTenSec:  requestsInTenSec,
		ResponsesInTenSec: responsesInTenSec,
	}
}

// NewEmptyRecovery allocates and returns a new empty Recovery.
func NewEmptyRecovery() *Recovery {
	return new(Recovery)
}

// Unmarshal parses the JSON-encoded recovery data.
func (r *Recovery) Unmarshal(data []byte) error {
	savedData := &Recovery{}

	err := json.Unmarshal(data, savedData)

	if nil != err {
		return err
	}

	r.Index = savedData.Index
	r.TimeStamp = savedData.TimeStamp
	r.RequestsInTenSec = savedData.RequestsInTenSec
	r.ResponsesInTenSec = savedData.ResponsesInTenSec

	return nil
}

// Marshal returns the JSON encoding of the recovery object.
func (r *Recovery) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

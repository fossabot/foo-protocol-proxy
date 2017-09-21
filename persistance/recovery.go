package persistance

import (
	"encoding/json"
)

type (
	Recovery struct {
		Index             uint8      `json:"index"`
		TimeStamp         uint64     `json:"time_stamp"`
		RequestsInTenSec  [10]uint64 `json:"requests_10s"`
		ResponsesInTenSec [10]uint64 `json:"responses_10s"`
	}
)

func NewRecovery(index uint8, timestamp uint64, requestsInTenSec, responsesInTenSec [10]uint64) *Recovery {
	return &Recovery{
		Index:             index,
		TimeStamp:         timestamp,
		RequestsInTenSec:  requestsInTenSec,
		ResponsesInTenSec: responsesInTenSec,
	}
}

func NewEmptyRecovery() *Recovery {
	return new(Recovery)
}

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

func (r *Recovery) Marshall() ([]byte, error) {
	return json.Marshal(r)
}

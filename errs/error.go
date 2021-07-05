package errs

import (
	"encoding/json"
	"errors"
	"github.com/PandaTtttt/go-assembly/util/must"
)

type Error struct {
	RetCode RetCode `json:"retCode"`
	RetMsg  string  `json:"retMsg"`
}

func (e *Error) Error() string {
	if e.RetMsg != "" {
		return e.RetCode.String() + "[" + e.RetMsg + "]"
	}
	return e.RetCode.String()
}

func (e *Error) Json() string {
	return string(must.Byte(json.Marshal(e)))
}

// Decode tries to decode the given bytes to errs.Err,
// returns a standard error instead if unmarshal failed.
func Decode(body []byte) error {
	var e *Error
	err := json.Unmarshal(body, &e)
	if err != nil {
		return errors.New(string(body))
	}
	return e
}

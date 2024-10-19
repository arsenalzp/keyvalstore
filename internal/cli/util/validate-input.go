// Package is used for validation of importing/exporting data.

package util

import (
	"encoding/json"
	"fmt"

	"github.com/arsenalzp/keyvalstore/internal/cli/errors"
)

const (
	KEY_LENGTH   = 256
	VALUE_LENGTH = 511
)

// simple structure is used to unmarshal incoming JSON data into it
// to check its validity
type dataStruct struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// function ValidateData is used to validate incoming JSON data;
// it shows error in case of JSON data is invalid, however
// it's possible to use json.Valid() function insted.
func ValidateData(data []byte) error {
	var dataStruct []dataStruct
	err := json.Unmarshal(data, &dataStruct)
	if err != nil {
		return err
	}

	return nil
}

// validate input of the key and the value parameters
func ValidateInput(key, value []byte) error {
	switch {
	case len(key) == 0:
		return errors.New("input validation error: key and value shoudn't be empty", errors.KeyEmptyErr, nil)
	case len(key) > KEY_LENGTH:
		message := fmt.Sprintf("input validation error: key size is greater than %d bytes, current size: %d", KEY_LENGTH, len(key))
		return errors.New(message, errors.KeyLenExceededErr, nil)
	case len(value) > VALUE_LENGTH:
		message := fmt.Sprintf("input validation error: value size is greater than %d bytes, current size: %d", VALUE_LENGTH, len(key))
		return errors.New(message, errors.ValueLenExceededErr, nil)
	default:
		return nil
	}
}

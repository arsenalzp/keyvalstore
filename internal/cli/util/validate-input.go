// Package is used for validation of importing/exporting data.

package util

import (
	"encoding/json"
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

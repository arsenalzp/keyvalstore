// Wrapper for standard errors package.
// It defines error codes and implements stack wrapper.

package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	InvalidAddrErr      = "ECLI-0001"
	NetworkErr          = "ECLI-0002"
	ReadStdinErr        = "ECLI-0003"
	ReadServerErr       = "ECLI-0004"
	WriteServerErr      = "ECLI-0005"
	InvalidExport       = "ECLI-0006"
	GetResponseError    = "ECLI-0007"
	SetResponseError    = "ECLI-0008"
	DelResponseError    = "ECLI-0009"
	ExpResponseError    = "ECLI-0010"
	ImpResponseError    = "ECLI-0011"
	ServerResponseError = 'N'
	KeyLenExceededErr   = "ECLI-0012"
	ValueLenExceededErr = "ECLI-1013"
	KeyEmptyErr         = "ECLI-2014"
)

type errorCmd struct {
	Err  error
	Code string
	Msg  string
}

func (e errorCmd) Error() string {
	// if error argument is nil - don't print it out
	if e.Err == nil {
		return fmt.Sprintf(
			"%s, code: %s",
			e.Msg, e.Code,
		)
	}

	return fmt.Sprintf(
		"%s, code: %s, %s",
		e.Msg, e.Code, e.Err,
	)
}

func New(msg string, code string, err error) error {
	return errors.WithStack(&errorCmd{
		Msg: msg, Code: code, Err: err,
	})
}

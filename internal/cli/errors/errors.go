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
)

type errorCmd struct {
	Err  error
	Code string
	Msg  string
}

func (e errorCmd) Error() string {
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

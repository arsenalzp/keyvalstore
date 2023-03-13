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
	WriteServerErr      = "ECLI-1005"
	InvalidExport       = "ECLI-0006"
	KeyEmptyErr         = "ECLI-0007"
	ValueEmptyErr       = "ECLI-1008"
	DelServerRespErr    = "ECLI-0009"
	GetServerRespErr    = "ECLI-1010"
	SetServerRespErr    = "ECLI-2011"
	ExpServerRespErr    = "ECLI-3012"
	ImpServerRespErr    = "ECLI-4013"
	DelCancelErr        = "ECLI-0014"
	GetCancelErr        = "ECLI-1015"
	SetCancelErr        = "ECLI-2016"
	ExpCancelErr        = "ECLI-3017"
	ImpCancelErr        = "ECLI-4018"
	ServerResponseError = 'N'
	InputValidationErr  = "ECLI-0019"
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

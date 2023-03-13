// Wrapper for standard errors package.
// It defines error codes and implements stack wrapper.

package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	NetworkErrIpaddrs = "ESRV-0001"
	NetworkErr        = "ESRV-1002"
	NetworkCallErr    = "ESRV-2003"
	ReadClientErr     = "ESRV-0004"
	WriteClientErr    = "ESRV-0005"
	ServerIntErr      = "ESRV-0006"
	UnknownClientOps  = "WSRV-0007"
	OperationTimeout  = "WSRV-0008"
	SettOpErr         = "ESRV-0009"
	GetOpErr          = "ESRV-0010"
	DelOpErr          = "ESRV-0011"
	ExpOpErr          = "ESRV-0012"
	ImpOpErr          = "ESRV-0013"
	SetOpTimeout      = "WSRV-0015"
	GetOpTimeout      = "WSRV-0016"
	DelOpTimeout      = "WSRV-0017"
	ExpOpTimeout      = "WSRV-0018"
	ImpOpTimeout      = "WSRV-0019"
	ExportErr         = "EKV-0020"
	ImportErr         = "EKV-0021"
	GetErr            = "EKV-0022"
	SetErr            = "EKV-0023"
	DelErr            = "EKV-0024"
	CRLExpiredErr     = "ETLS-0025"
	CRLLoadErr        = "ETLS-1026"
	CRLParseErr       = "ETLS-2027"
	KeyCertLoadErr    = "ETLS-0028"
	CAcertLoadErr     = "ETLS-0029"
	CAPoolLoadErr     = "ETLS-1030"
	CRLValidErr       = "ETLS-3031"
	CRLCertRevokErr   = "ETLS-4032"
	StorageNilErr     = "ESTRG-0033"
	StorageInitErr    = "ESTRG-1034"
	StorageKindErr    = "ESTRG-2035"
	StorageKindUndef  = "ESTRG-3036"
	CRLOpenErr        = "ETLS-3037"
	CRLStatErr        = "ETLS-4038"
	HashTabInsErr     = "EHTAB-0039"
	HashTabDelErr     = "EHTAB-1040"
	HashTabExpErr     = "EHTAB-2041"
	HashTabImpErr     = "EHTAB-3042"
	HashTabSrchErr    = "EHTAB-4043"
	NetworkInitErr    = "ESRV-3044"
	NetworkInitTLSErr = "ESRV-4045"
	SrvStartErr       = "ESRV-5046"
	SrvStopErr        = "ESRV-6047"
)

type errCommon struct {
	Msg  string
	Code string
	Err  error
}

func (e errCommon) Error() string {
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
	return errors.WithStack(&errCommon{
		Msg: msg, Code: code, Err: err,
	})
}

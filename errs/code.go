package errs

import (
	"fmt"
	"strconv"
)

type RetCode uint32

const (
	// Internal is the generic error that maps to HTTP 500.
	Internal RetCode = iota + 100001
	// NotFound indicates a given resource is not found.
	NotFound
	// Forbidden indicates the user doesn't have the permission to
	// perform given operation.
	Forbidden
	// Unauthenticated indicates the oauth2 authentication failed.
	Unauthenticated
	// InvalidArgument indicates the input is invalid.
	InvalidArgument
	// InvalidConfig indicates the config is invalid.
	InvalidConfig
	// Conflict indicates a database transactional conflict happens.
	Conflict
	// TryAgain indicates a temporary outage and retry
	// could eventually lead to success.
	TryAgain
)

func (c RetCode) String() string {
	switch c {
	case Internal:
		return "Internal"
	case NotFound:
		return "NotFound"
	case Forbidden:
		return "Forbidden"
	case Unauthenticated:
		return "Unauthenticated"
	case InvalidArgument:
		return "InvalidArgument"
	case InvalidConfig:
		return "InvalidConfig"
	case Conflict:
		return "Conflict"
	case TryAgain:
		return "TryAgain"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

func (c RetCode) Wrap(err error) error {
	_, ok := err.(*Error)
	if ok {
		return err
	}
	
	return &Error{
		RetCode: c,
		RetMsg:  err.Error(),
	}
}

func (c RetCode) New(a ...string) error {
	msg := ""
	for i, s := range a {
		if i > 0 {
			msg += " "
		}
		msg += s
	}
	return &Error{
		RetCode: c,
		RetMsg:  msg,
	}
}

func (c RetCode) Newf(msg string, args ...interface{}) error {
	return &Error{
		RetCode: c,
		RetMsg:  fmt.Sprintf(msg, args...),
	}
}

func (c RetCode) Is(err error) bool {
	v, ok := err.(*Error)
	if !ok {
		// all other errors are internal error.
		return c == Internal
	}
	return v.RetCode == c
}

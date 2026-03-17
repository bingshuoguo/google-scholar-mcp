package scholar

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrNoResults           = errors.New("no results")
	ErrTimeout             = errors.New("request timeout")
	ErrUpstreamBlocked     = errors.New("upstream blocked")
	ErrUpstreamUnavailable = errors.New("upstream unavailable")
	ErrParseFailed         = errors.New("parse failed")
)

type Error struct {
	Kind    error
	Message string
	Cause   error
}

func (e *Error) Error() string {
	switch {
	case e.Message != "" && e.Cause != nil:
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	case e.Message != "":
		return e.Message
	case e.Cause != nil:
		return e.Cause.Error()
	default:
		return "unknown scholar error"
	}
}

func (e *Error) Unwrap() error {
	if e.Kind != nil {
		return e.Kind
	}
	return e.Cause
}

func wrap(kind error, message string, cause error) error {
	return &Error{Kind: kind, Message: message, Cause: cause}
}

func ClassifyMessage(err error) string {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return err.Error()
	case errors.Is(err, ErrNoResults):
		return err.Error()
	case errors.Is(err, ErrTimeout):
		return "Google Scholar request timed out. Try again with a shorter query or later."
	case errors.Is(err, ErrUpstreamBlocked):
		return "Google Scholar blocked the request. Reduce request frequency or try again later."
	case errors.Is(err, ErrUpstreamUnavailable):
		return "Google Scholar is temporarily unavailable. Try again later."
	case errors.Is(err, ErrParseFailed):
		return "Google Scholar returned a page that could not be parsed. The page structure may have changed."
	default:
		return "Unexpected Google Scholar error. Check server logs for details."
	}
}

package httpx

import (
	"net/http"
	"strings"
)

// HTTP Errors
var (
	ErrUnsupportedMediaType        = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound                    = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                   = NewHTTPError(http.StatusForbidden)
	ErrMethodNotAllowed            = NewHTTPError(http.StatusMethodNotAllowed)
	ErrStatusRequestEntityTooLarge = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrTooManyRequests             = NewHTTPError(http.StatusTooManyRequests)
	ErrBadRequest                  = NewHTTPError(http.StatusBadRequest)
	ErrBadGateway                  = NewHTTPError(http.StatusBadGateway)
	ErrInternalServerError         = NewHTTPError(http.StatusInternalServerError)
	ErrRequestTimeout              = NewHTTPError(http.StatusRequestTimeout)
	ErrServiceUnavailable          = NewHTTPError(http.StatusServiceUnavailable)
)

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Err     error
	Code    int
	Message string
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message ...string) *HTTPError {
	e := &HTTPError{Code: code, Message: strings.ToLower(http.StatusText(code))}
	if len(message) > 0 {
		e.Message = strings.Join(message, ": ")
	}
	return e
}

// WrapHTTPError creates a new HTTPError instance with internal error set.
func WrapHTTPError(err error, code int, message ...string) *HTTPError {
	e := NewHTTPError(code, message...)
	e.Err = err
	return e
}

// Unwrap satisfies the Go 1.13 error wrapper interface.
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// Error makes it compatible with error interface.
func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// WithError returns the same HTTPError instance with err set to HTTPError.err field
func (e *HTTPError) WithError(err error) *HTTPError {
	e.Err = err
	return e
}

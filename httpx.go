package httpx

import (
	"errors"
	"log"
	"net/http"
	"os"
)

var (
	// Logger is the global logger used to log library-specific internal errors.
	// Set this global variable to your preferred logger. Defaults to stdlib log.
	Logger logger = log.New(os.Stderr, "httpx: ", log.LstdFlags|log.Lmsgprefix)

	// HTTPErrorHandler is used to handle all HTTP errors returned by every HTTP handler.
	// Set this global variable to customise the behaviour.
	HTTPErrorHandler HTTPErrorHandlerFunc = HandleHTTPError(false)
)

type logger interface{ Println(v ...any) }

// HTTPErrorHandlerFunc is a centralised handler for HTTPError.
type HTTPErrorHandlerFunc func(req *Request, res *Responder, err error)

// HandleHTTPError returns the default HTTPErrorHandler used.
// If expose is true, returned response will be the internal error message.
func HandleHTTPError(expose bool) HTTPErrorHandlerFunc {
	return func(req *Request, res *Responder, err error) {
		if res.Committed {
			return
		}

		e := &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
		errors.As(err, &e)
		res.Status(e.Code)

		if expose {
			e.Message = err.Error()
		}

		var resErr error
		if req.Method == http.MethodHead {
			resErr = res.NoContent()
		} else {
			resErr = res.String(e.Message)
		}

		if resErr != nil {
			Logger.Println(resErr) // rare error case
		}
	}
}

type contextKey string

var requestKey = contextKey("request")

// HandlerFunc is an adapter to allow the use of ordinary functions as HTTP handlers,
// with *Request and *Responder as parameters.
//
// If f is a function with the appropriate signature, HandlerFunc(f) is a http.Handler that calls f.
type HandlerFunc func(req *Request, res *Responder) error

// ServeHTTP wraps http.Request into Request and http.ResponseWriter into Responder
// before passing them into and call h(req, res).
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, ok := r.Context().Value(requestKey).(*Request)
	if !ok {
		req = NewRequest(r)
		req.SetValue(requestKey, req)
	}
	res, ok := w.(*Responder)
	if !ok {
		res = NewResponder(w)
	}
	if err := h(req, res); err != nil {
		HTTPErrorHandler(req, res, err)
	}
}

// H is a convenient adapter that wraps the translation of http.Handler to HandlerFunc.
func H(handler http.Handler) HandlerFunc {
	return func(req *Request, res *Responder) error {
		handler.ServeHTTP(res, req.Request)
		return nil
	}
}

package httpx

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const defaultMaxMemory = 32 << 20 // 32 MB

// RequestBinder is used to bind request body data to objects.
// Set this global variable to your preferred binder once before
// calling any Request.Bind. The default DefaultRequestBinder binds data
// according to the "Content-Type" header.
var RequestBinder interface {
	Bind(req *Request, v any) error
} = new(DefaultRequestBinder)

// Request wraps an *http.Request.
// See: https://golang.org/pkg/net/http/#Request
type Request struct {
	*http.Request // inherit from http.Request

	query url.Values
}

// NewRequest creates a new instance of Request.
func NewRequest(r *http.Request) *Request {
	return &Request{Request: r}
}

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (r *Request) IsTLS() bool {
	return r.Request.TLS != nil
}

// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
func (r *Request) IsWebSocket() bool {
	upgrade := r.Request.Header.Get(HeaderUpgrade)
	return strings.EqualFold(upgrade, "websocket")
}

// Scheme returns the HTTP protocol scheme, http or https.
func (r *Request) Scheme() string {
	// Can't use Request.URL.Scheme
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if r.IsTLS() {
		return "https"
	}
	if scheme := r.Request.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := r.Request.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := r.Request.Header.Get(HeaderXForwardedSSL); ssl == "on" {
		return "https"
	}
	if scheme := r.Request.Header.Get(HeaderXURLScheme); scheme != "" {
		return scheme
	}
	return "http"
}

// QueryParam returns the query param for the provided name.
func (r *Request) QueryParam(name string) string {
	if r.query == nil {
		r.query = r.Request.URL.Query()
	}
	return r.query.Get(name)
}

// QueryParams returns the query parameters as url.Values.
func (r *Request) QueryParams() url.Values {
	if r.query == nil {
		r.query = r.Request.URL.Query()
	}
	return r.query
}

// QueryString returns the URL query string.
func (r *Request) QueryString() string {
	return r.Request.URL.RawQuery
}

// FormParams returns the form parameters as url.Values.
func (r *Request) FormParams() (url.Values, error) {
	if strings.HasPrefix(r.Header.Get(HeaderContentType), MIMEMultipartForm) {
		if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
	}
	return r.Form, nil
}

// FormFile returns the multipart form file for the provided name.
func (r *Request) FormFile(name string) (*multipart.FileHeader, error) {
	f, fh, err := r.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil
}

// MultipartForm returns the multipart form.
func (r *Request) MultipartForm() (*multipart.Form, error) {
	err := r.Request.ParseMultipartForm(defaultMaxMemory)
	return r.Request.MultipartForm, err
}

// SetValue sets a value with key to the underlying http.Request's context.Context.
// The context can be retrieved using Request.Context().
func (r *Request) SetValue(key, val any) {
	ctx := context.WithValue(r.Request.Context(), key, val)
	r.Request = r.Request.WithContext(ctx)
}

// GetValue gets a value by key from the underlying http.Request's context.Context.
// The context can be retrieved using Request.Context().
func (r *Request) GetValue(key any) any {
	return r.Request.Context().Value(key)
}

// Bind binds data from request body to v.
// Immediately panic if RequestBinder is not set in advance.
func (r *Request) Bind(v any) error {
	if RequestBinder == nil {
		panic("undefined request binder")
	}
	return RequestBinder.Bind(r, v)
}

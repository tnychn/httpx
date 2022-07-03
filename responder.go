package httpx

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"errors"
	"html/template"
	"io"
	"net"
	"net/http"
)

var (
	errInvalidRedirectCode    = errors.New("invalid redirect status code")
	errHeaderAlreadyCommitted = errors.New("response already committed")
)

// Responder wraps an http.ResponseWriter and implements its interface to be used
// by an HTTP handler to construct an HTTP response.
// See: https://golang.org/pkg/net/http/#ResponseWriter
type Responder struct {
	beforeFuncs []func()
	afterFuncs  []func()

	Size       int64
	Committed  bool
	StatusCode int

	Writer http.ResponseWriter
}

// NewResponder creates a new instance of Responder.
func NewResponder(w http.ResponseWriter) *Responder {
	return &Responder{Writer: w}
}

func (r *Responder) writeContentType(value string) {
	header := r.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

func (r *Responder) writeHeader() {
	if r.StatusCode <= 0 {
		r.StatusCode = http.StatusOK
	}
	r.WriteHeader(r.StatusCode)
}

// Header returns the header map for the writer that will be sent by
// WriteHeader. Changing the header after a call to WriteHeader (or Write) has
// no effect unless the modified headers were declared as trailers by setting
// the "Trailer" header before the call to WriteHeader (see example)
// To suppress implicit response headers, set their value to nil.
// Example: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
func (r *Responder) Header() http.Header {
	return r.Writer.Header()
}

// Before registers a function which is called just before the response is written.
func (r *Responder) Before(fn func()) {
	r.beforeFuncs = append(r.beforeFuncs, fn)
}

// After registers a function which is called just after the response is written.
// If the Content-Length is unknown, none of the after function is executed.
func (r *Responder) After(fn func()) {
	r.afterFuncs = append(r.afterFuncs, fn)
}

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *Responder) WriteHeader(code int) {
	if r.Committed {
		Logger.Println(errHeaderAlreadyCommitted)
		return
	}
	r.StatusCode = code
	for _, fn := range r.beforeFuncs {
		fn()
	}
	r.Writer.WriteHeader(r.StatusCode)
	r.Committed = true
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *Responder) Write(b []byte) (n int, err error) {
	if !r.Committed {
		r.writeHeader()
	}
	n, err = r.Writer.Write(b)
	r.Size += int64(n)
	for _, fn := range r.afterFuncs {
		fn()
	}
	return
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Responder) Flush() {
	r.Writer.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (r *Responder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.Writer.(http.Hijacker).Hijack()
}

// SetCookie adds a Set-Cookie header in HTTP response.
func (r *Responder) SetCookie(cookie *http.Cookie) *Responder {
	http.SetCookie(r, cookie)
	return r
}

// Status sets the status code of the response without committing it.
func (r *Responder) Status(code int) *Responder {
	r.StatusCode = code
	return r
}

// HTML sends an HTTP response.
func (r *Responder) HTML(html string) error {
	return r.Blob(MIMETextHTMLCharsetUTF8, []byte(html))
}

// String sends a string response.
func (r *Responder) String(s string) error {
	return r.Blob(MIMETextPlainCharsetUTF8, []byte(s))
}

// JSON sends a JSON response.
func (r *Responder) JSON(i any, indent string) error {
	r.writeContentType(MIMEApplicationJSONCharsetUTF8)
	r.writeHeader()
	enc := json.NewEncoder(r)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// XML sends an XML response.
func (r *Responder) XML(i any, indent string) error {
	r.writeContentType(MIMEApplicationXMLCharsetUTF8)
	r.writeHeader()
	enc := xml.NewEncoder(r)
	if indent != "" {
		enc.Indent("", indent)
	}
	if _, err := r.Write([]byte(xml.Header)); err != nil {
		return err
	}
	return enc.Encode(i)
}

// Blob sends a blob response with content type.
func (r *Responder) Blob(contentType string, b []byte) error {
	r.writeContentType(contentType)
	r.writeHeader()
	_, err := r.Write(b)
	return err
}

// Stream sends a streaming response with content type.
func (r *Responder) Stream(contentType string, reader io.Reader) error {
	r.writeContentType(contentType)
	r.writeHeader()
	_, err := io.Copy(r, reader)
	return err
}

// Template renders a template with data and sends a text/html response.
func (r *Responder) Template(tpl *template.Template, name string, data any) error {
	r.writeContentType(MIMETextHTMLCharsetUTF8)
	r.writeHeader()
	return tpl.ExecuteTemplate(r, name, data)
}

// NoContent sends a response with no body.
func (r *Responder) NoContent() error {
	r.writeHeader()
	return nil
}

// Redirect redirects the request to a provided URL.
func (r *Responder) Redirect(url string) error {
	if r.StatusCode < 300 || r.StatusCode > 308 {
		return errInvalidRedirectCode
	}
	r.Header().Set(HeaderLocation, url)
	r.writeHeader()
	return nil
}

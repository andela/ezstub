package main

import "net/http"

// responseWriter wraps http.ResponseWriter to extract
// information for log.
type responseWriter struct {
	rw            http.ResponseWriter
	headerWritten bool
	status        int
	written       int
}

// Header returns the header map that will be sent by
// WriteHeader. Changing the header after a call to
// WriteHeader (or Write) has no effect unless the modified
// headers were declared as trailers by setting the
// "Trailer" header before the call to WriteHeader (see example).
// To suppress implicit response headers, set their value to nil.
func (r *responseWriter) Header() http.Header {
	return r.rw.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (r *responseWriter) Write(b []byte) (int, error) {
	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.rw.Write(b)
	r.written += n
	return n, err
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (r *responseWriter) WriteHeader(status int) {
	r.status = status
	r.headerWritten = true
	r.rw.WriteHeader(status)
}

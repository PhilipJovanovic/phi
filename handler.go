package phi

import (
	"errors"
	"net/http"
)

var (
	defaultHandler = func(w http.ResponseWriter, r *http.Request, e *Error) {
		w.Write([]byte("unknown error"))
	}
)

type Handler func(w *Response, r *Request) *Error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(
		&Response{
			ResponseWriter: w,
		},
		&Request{
			Request: r,
		},
	); err != nil {
		defaultHandler(w, r, err)
	}
}

// set custom error handler
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, e *Error)

func SetErrorHandler(fn ErrorHandlerFunc) error {
	if fn == nil {
		return errors.New("couldn't set empty error handling function")
	}

	defaultHandler = fn
	return nil
}

package phi

import (
	"errors"
	"net/http"
)

var (
	handler = func(w http.ResponseWriter, r *http.Request, e *Error) {
		w.Write([]byte("unknown error"))
	}
)

type Error struct {
	Error      string
	Message    string
	StatusCode int
}

type ErrorHandler func(w http.ResponseWriter, r *http.Request) *Error

func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		handler(w, r, err)
	}
}

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, e *Error)

func SetErrorHandler(fn ErrorHandlerFunc) error {
	if fn == nil {
		return errors.New("couldn't set empty error handling function")
	}

	handler = fn
	return nil
}

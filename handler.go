package phi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	defaultHandler = func(w http.ResponseWriter, r *http.Request, e *Error) {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}

		parsed, err := json.Marshal(Map{
			"error":   e.Error,
			"message": e.Message,
		})
		if err == nil {
			if _, err = w.Write(parsed); err != nil {
				log.Printf("#> defaultHandler: %v", err)
			}

			return
		}

		if e.StatusCode != 0 {
			w.WriteHeader(e.StatusCode)
		}

		if _, err = w.Write(parsed); err != nil {
			log.Printf("#> defaultHandler: %v", err)
		}
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

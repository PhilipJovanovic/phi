package phi

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func (w Response) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w Response) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w Response) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

// send response with application/json
func (w Response) JSON(data interface{}) *Error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	parsed, err := json.Marshal(Map{
		"data": data,
	})
	if err != nil {
		return &parseError
	}

	if _, er := w.Write(parsed); er != nil {
		return &writingError
	}

	return nil
}

// send response with contentType
func (w Response) Response(data []byte, contentType string) *Error {
	w.Header().Set("Content-Type", contentType)

	if _, err := w.Write(data); err != nil {
		return &writingError
	}

	return nil
}

// send error response with status code
func (w Response) Error(err *Error) *Error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	parsed, er := json.Marshal(Map{
		"error":   err.Error,
		"message": err.Message,
	})
	if er != nil {
		return &parseError
	}

	if err.StatusCode != 0 {
		w.WriteHeader(err.StatusCode)
	}

	if _, er = w.Write(parsed); er != nil {
		return &writingError
	}

	return nil
}

// send redirect response
func (w Response) Redirect(req *http.Request, location string, code int) *Error {
	http.Redirect(w.ResponseWriter, req, location, code)
	return nil
}

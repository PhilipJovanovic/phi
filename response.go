package phi

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

// send response with application/json
//
// content will be wrapped in { "data" : <content> } object
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
//
// default error will look like this:
//
//	{
//	  "error": "unknownError",
//	  "message": err.Error()
//	}
func (w Response) Error(err error) *Error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	parsed, er := json.Marshal(Map{
		"error":   "unknownError",
		"message": err.Error(),
	})
	if er != nil {
		return &parseError
	}

	if _, er = w.Write(parsed); er != nil {
		return &writingError
	}

	return nil
}

// send error response with custom status code
func (w Response) ErrorCustomStatus(err error, statusCode int) *Error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	parsed, er := json.Marshal(Map{
		"error":   "unknownError",
		"message": err.Error(),
	})
	if er != nil {
		return &parseError
	}

	w.WriteHeader(statusCode)

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

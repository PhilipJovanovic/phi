package phi

var (
	writingError = Error{
		Error:   "writingError",
		Message: "error while writing response",
	}

	parseError = Error{
		Error:   "parseError",
		Message: "error while parsing response",
	}

	decodingError = Error{
		Error:   "decodingError",
		Message: "error while decoding request body",
	}
)

type Error struct {
	Error      string
	Message    string
	StatusCode int
}

// Validation error can be used for validating post bodies
func ValidatingError(e error) *Error {
	return &Error{
		Error:      "validatingError",
		Message:    e.Error(),
		StatusCode: 400,
	}
}

// Parameter error for error handling regarding parameter missing /yeet/{cid} -> cid = parameter
func URLParameterError(e string) *Error {
	return &Error{
		Error:      "missingURLParameters",
		Message:    e,
		StatusCode: 400,
	}
}

// Query Parameter error for error handling regarding parameter missing /yeet?id=1337 -> id = query parameter
func QueryParameterError(e string) *Error {
	return &Error{
		Error:      "missingQueryParameters",
		Message:    e,
		StatusCode: 400,
	}
}

// Body Parameter error for error handling regarding parameter missing POST body = { "data": "123"} -> data = body parameter
func BodyParameterError(e string) *Error {
	return &Error{
		Error:      "missingBodyParameters",
		Message:    e,
		StatusCode: 400,
	}
}

// Unknown error for generic error handling
func UnknownError(e error) *Error {
	return &Error{
		Error:      "unknownError",
		Message:    e.Error(),
		StatusCode: 500,
	}
}

// Unauthorized error + statuscode 401
func Unauthorized() *Error {
	return &Error{
		Error:      "unauthorized",
		Message:    "invalid token",
		StatusCode: 401,
	}
}

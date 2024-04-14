package phi

var (
	writingError = Error{
		Error:   "writingError",
		Message: "error writing response",
	}
)

// Validation error can be used for validating post bodies
func ValidatingError(e error) *Error {
	return &Error{
		Error:   "validatingError",
		Message: e.Error(),
	}
}

// Parameter error for error handling regarding parameter missing /yeet/{cid} -> cid = parameter
func ParameterError(e string) *Error {
	return &Error{
		Error:   "parameter(s) missing",
		Message: e,
	}
}

// Unknown error for generic error handling
func UnknownError(e error) *Error {
	return &Error{
		Error:   "unknownError",
		Message: e.Error(),
	}
}

// Unauthorized error
func Unauthorized() *Error {
	return &Error{
		Error:   "unauthorized",
		Message: "invalid token",
	}
}

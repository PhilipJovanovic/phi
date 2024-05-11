package phi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type Request struct {
	*http.Request
}

// /yeet/{cid} -> cid = parameter
func (r *Request) URLParam(id string) (string, *Error) {
	if s := URLParam(r.Request, id); s != "" {
		return s, nil
	}

	return "", URLParameterError(id)
}

// /yeet?id=1337 -> id = query parameter
func (r *Request) QueryParam(id string) (string, *Error) {
	if s := r.FormValue(id); s != "" {
		return s, nil
	}

	return "", QueryParameterError(id)
}

// if JWTAuth middleware is used, this function can be used to extract the token from the context
/* func (r *Request) GetContextToken(id string) *middleware.Token {
	token := r.Context().Value(middleware.TOKEN_CONTEXT).(middleware.Token)

	return &token
} */

// Validate post bodies
// Example:
//
//	type Body struct {
//		Data string `json:"data,required"` 	// required
//		Dutu string `json:"dutu"`		// optional
//	}
func Validate[T any](r *Request) (*T, *Error) {
	var body T

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, &decodingError
	}

	fields := reflect.ValueOf(&body).Elem()
	errorList := make([]string, 0)

	for i := 0; i < fields.NumField(); i++ {
		jsonTags := fields.Type().Field(i).Tag.Get("json")

		if strings.Contains(jsonTags, "required") && fields.Field(i).IsZero() {
			errorList = append(errorList, strings.Split(jsonTags, ",")[0])
		}
	}

	if len(errorList) > 0 {
		return nil, BodyParameterError(fmt.Sprintf("missing '%s'", strings.Join(errorList, ",")))
	}

	return &body, nil
}

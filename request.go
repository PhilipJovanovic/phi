package phi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// A Wrapper for
//
//	ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
//	r.WithContext(ctx)
//
// mostly used by middlewares f.E.:
//
//	func TestMW(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			str := "test"
//			ctx := context.WithValue(r.Context(), "CONTEXT_ID", &str)
//			next.ServeHTTP(w, r.WithContext(ctx))
//		})
//	}
func (r *Request) SetContext(contextId string, data any) *http.Request {
	ctx := context.WithValue(r.Context(), contextId, data)

	return r.WithContext(ctx)
}

// A Wrapper for Context().Value to return a casted value
func GetContext[T any](r *Request, contextId string) *T {
	t := r.Context().Value(contextId).(*T)

	return t
}

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

// Validate post bodies but accepts string
// Example:
//
//	type Body struct {
//		Data string `json:"data,required"` 	// required
//		Dutu string `json:"dutu"`		// optional
//	}
func ValidateString[T any](r string) (*T, *Error) {
	var body T

	reader := io.NopCloser(strings.NewReader(r))

	if err := json.NewDecoder(reader).Decode(&body); err != nil {
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

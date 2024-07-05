package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.philip.id/phi"
	"go.philip.id/phi/jwtauth"
)

type Token struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
}

type TOKEN_TYPE string

// Use this type to get/set the token from the context
const TOKEN_CONTEXT TOKEN_TYPE = "_token"

var (
	unauthorizedFunc = phi.Unauthorized
	tokenCheckFunc   = ImplementAccessCheck

	invalidUserId = phi.Error{
		Error:      "invalidUserId",
		Message:    "user id is not primitive.ObjectID",
		StatusCode: 400,
	}

	tokenNotFound = phi.Error{
		Error:      "tokenNotFound",
		Message:    "context token could not be found",
		StatusCode: 400,
	}
)

// SetUnauthorizedFunc sets the function to be called when a request is unauthorized
//
// default is phi.Unauthorized
func SetUnauthorizedFunc(fn func() *phi.Error) {
	unauthorizedFunc = fn
}

// set new tokencheck function, f.e.. check username, password against a database
//
// Example (mongopiet):
//
//	func TokenCheck(username, password string) (*Token, error) {
//		token, err := database.FindOne("apiTokens", bson.M{"token": username})
//		if err != nil {
//			return nil, errors.New("not found")
//		}
//
//		return &phi.Token{
//			ID:     a.User.Hex(),
//			Subject:  a.Subject,
//		}
//	}
func SetTokenCheckFunc(fn func(username, password string) (*Token, error)) {
	tokenCheckFunc = fn
}

// default implementation of token check, wont work in production!
func ImplementAccessCheck(username, password string) (*Token, error) {
	return nil, errors.New("not implemented")
}

// GetToken returns the token from the context
func GetToken(r *phi.Request) *Token {
	if token, ok := r.Context().Value(TOKEN_CONTEXT).(Token); ok {
		return &token
	}

	return nil
}

// Opinionated helper function to get the user id from the token as primitive.ObjectID
func GetUserID(r *phi.Request) (*primitive.ObjectID, *phi.Error) {
	token, ok := r.Context().Value(TOKEN_CONTEXT).(Token)
	if !ok {
		return nil, &tokenNotFound
	}

	id, err := primitive.ObjectIDFromHex(token.ID)
	if err != nil {
		return nil, &invalidUserId
	}

	return &id, nil
}

// Checks for bearer token or basic auth and returns unauthorized if not found
//
// unauthorized response can be set via SetTokenCheckFunc
//
// Can be used for endpoints which are gonna be used for a frontend and from an api
// at the same time
//
// Tokens can be extracted like one of the following:
//
//	token := r.Context().Value(middleware.TOKEN_CONTEXT).(middleware.Token)
//	token := phi.GetToken(r) // only works with *phi.Request
func JWTOrAPIAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, err := checkBearer(r); err == nil {
			ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if token, err := checkBasic(r); err == nil {
			ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		phi.ErrorHandler(w, r, unauthorizedFunc())
	})
}

// Same as JWTOrAPIAuth but continues without adding the token if unauthorized
//
// Can be used for cases where an authenticated user will receive a different
// response but still has access to the ressource
func JWTOrAPIAuthOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, err := checkBearer(r); err == nil {
			ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if token, err := checkBasic(r); err == nil {
			ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Checks for bearer token and returns unauthorized if not found
//
// unauthorized response can be set via SetTokenCheckFunc
//
// Can be used for frontend authentication, middleware expects jwt token at every request,
// needs to be refreshed after expiry
//
// Tokens can be extracted like one of the following:
//
//	token := r.Context().Value(middleware.TOKEN_CONTEXT).(middleware.Token)
//	token := phi.GetToken(r) // only works with *phi.Request
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := checkBearer(r)
		if err != nil {
			phi.ErrorHandler(w, r, unauthorizedFunc())
			return
		}

		ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Same as JWTAuth but continues without adding the token if unauthorized
//
// Can be used for cases where an authenticated user will receive a different
// response but still has access to the ressource
func JWTAuthOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := checkBearer(r)
		if err == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Checks for basic authentication and returns unauthorized if not found
//
// unauthorized response can be set via SetTokenCheckFunc
//
// Can be used for api authentication, middleware expects basic header at every request,
// token is more likely to be longer available and should not be exposed to the client
//
// Tokens can be extracted like one of the following:
//
//	token := r.Context().Value(middleware.TOKEN_CONTEXT).(middleware.Token)
//	token := phi.GetToken(r) // only works with *phi.Request
func APIAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := checkBasic(r)
		if err != nil {
			phi.ErrorHandler(w, r, unauthorizedFunc())
			return
		}

		ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Same as APIAuth but continues without adding the token if unauthorized
//
// Can be used for cases where an authenticated user will receive a different
// response but still has access to the ressource
func APIAuthOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := checkBasic(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), TOKEN_CONTEXT, *token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JWT Bearer Token Validation
func checkBearer(r *http.Request) (*Token, error) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return nil, err
	}

	if token != nil && jwt.Validate(token) == nil {
		t := &Token{}

		// Get only ID for now
		t.ID = token.JwtID()
		t.Subject = token.Subject()

		// TODO: add more claim parsing here

		return t, nil
	}

	return nil, errors.New("token invalid")
}

// Basic Auth Validation
func checkBasic(r *http.Request) (*Token, error) {
	if username, password, ok := r.BasicAuth(); ok {
		t, err := tokenCheckFunc(username, password)
		if err != nil {
			return nil, err
		}

		return t, nil
	}

	return nil, errors.New("no basic auth found")
}

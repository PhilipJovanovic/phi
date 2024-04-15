package middleware

import (
	"context"
	"net/http"

	"github.com/lestrrat-go/jwx/jwt"
	"go.philip.id/phi/jwtauth"
)

// Checks for jwt bearer token
func BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// * JWT Bearer Token Validation
		if token := checkBearer(r); token != nil {
			ctx := context.WithValue(r.Context(), "userToken", *token)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	})
}

func checkBearer(r *http.Request) *auth.UserTokenObject {
	// * JWT Bearer Token Validation
	token, _, err := jwtauth.FromContext(r.Context())

	if err == nil && token != nil && jwt.Validate(token) == nil {
		if i, ok := token.Get("i"); ok {
			scope := Unknown

			for _, role := range token.PrivateClaims()["roles"].([]interface{}) {
				// There are only 3 types of bearer (only for the dash)
				if role.(string) == "ADMIN" {
					scope = Admin
				}

				if role.(string) == "USER" {
					scope = Regular
				}

				if role.(string) == "SUPPORT" {
					scope = Support
				}
			}

			return &auth.UserTokenObject{
				ID:     i.(string),
				Scope:  scope,
				Source: "_nofy_",
			}
		}
	}

	return nil
}

package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.philip.id/phi"
	"go.philip.id/phi/jwtauth"
)

var (
	r            *phi.Mux
	workingToken string
)

type CustomClaims struct {
	Source string `json:"source"`
	jwt.RegisteredClaims
}

func (m CustomClaims) Validate() error {
	if m.Source == "" {
		return errors.New("must be not empty")
	}
	return nil
}

func TestMain(t *testing.T) {
	r = phi.NewRouter()

	r.Use(jwtauth.Verifier(jwtauth.New("HS256", []byte("ultra_secret_key"), nil)))

	// No Auth at all
	r.Route("/noAuth", func(r phi.Router) {
		r.GET("/", handler)
	})

	r.Route("/bearer", func(r phi.Router) {
		r.Use(JWTAuth)
		r.GET("/", handler)
	})

	// sign a token for testing
	claims := &jwt.RegisteredClaims{
		ID:        "testToken",
		Subject:   "testLibrary",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}

	j := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jToken, err := j.SignedString([]byte("ultra_secret_key"))
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}
	workingToken = jToken

	t.Run("NoAuth Test", NoAuthMW)

	t.Run("Bearer Middleware Test", TestBearerMW)
	//t.Run("New Handler 2", TestNewHandler2)
}

func handler(res *phi.Response, req *phi.Request) *phi.Error {
	return res.JSON("success")
}

func NoAuthMW(t *testing.T) {
	req := httptest.NewRequest("GET", "/noAuth", nil)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	fmt.Printf("Body: %s\n", body)
}

func TestBearerMW(t *testing.T) {
	req := httptest.NewRequest("GET", "/bearer", nil)

	fmt.Println("Working Token: ", workingToken)
	req.Header.Set("Authorization", "Bearer "+workingToken)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	if string(body) != "{\"data\":\"success\"}" {
		t.Errorf("expected body to be {\"data\":\"success\"} got %s", body)
	}
}

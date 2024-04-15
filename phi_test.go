package phi

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var r *Mux

func TestMain(t *testing.T) {
	r = NewRouter()

	r.Route("/", func(r Router) {
		r.Get("/standard", handleRoot)
		r.GET("/newHandler", handleNewHandler)
	})

	t.Run("Test Standard", TestStandard)
	t.Run("New Handler", TestNewHandler)
}

func TestStandard(t *testing.T) {
	req := httptest.NewRequest("GET", "/standard", nil)

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

func TestErrorSuccess(t *testing.T) {
	req := httptest.NewRequest("GET", "/errorSuccess", nil)

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

func TestError(t *testing.T) {
	req := httptest.NewRequest("GET", "/error", nil)

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

func TestNewHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/newHandler", nil)

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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func handleNewHandler(res *Response, req *http.Request) *Error {
	return res.JSON("Hello World!")
}

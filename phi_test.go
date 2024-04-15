package phi

import (
	"bytes"
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
		r.POST("/newHandler", handleNewHandler)
	})

	t.Run("Test Standard", TestStandard)
	t.Run("New Handler", TestNewHandler)
	t.Run("New Handler 2", TestNewHandler2)
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
	data := bytes.NewBuffer([]byte("{\"data\":\"1337\"}"))
	req := httptest.NewRequest("POST", "/newHandler", data)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	expected := "{\"data\":\"success\"}"
	if fmt.Sprintf("%s", body) != expected {
		t.Errorf("expected body to be {\"data\":\"success\"} got %s", body)
	}
}

func TestNewHandler2(t *testing.T) {
	data := bytes.NewBuffer([]byte("{\"data\":\"1338\"}"))
	req := httptest.NewRequest("POST", "/newHandler", data)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	expected := "{\"error\":\"Invalid Data\",\"message\":\"Data is not 1337\"}"
	if fmt.Sprintf("%s", body) != expected {
		t.Errorf("expected body to be {\"data\":\"success\"} got %s", body)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func handleNewHandler(res *Response, req *Request) *Error {
	body, err := Validate[struct {
		Data string `json:"data,required"`
	}](req)
	if err != nil {
		return err
	}

	if body.Data != "1337" {
		return &Error{
			Error:   "Invalid Data",
			Message: "Data is not 1337",
		}
	}

	return res.JSON("success")
}

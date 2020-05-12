// +build unit_tests all_tests

package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankur-anand/prod-app/pkg/logger"
)

var logTesting, _ = logger.NewDevelopment()

func TestStaticHandler_Home(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	sh := newStaticHandler(logTesting)
	handler := http.HandlerFunc(sh.home)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"message": "hello world from rest svc"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestStaticHandler_HealthAliveness(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/health/live", nil)
	rr := httptest.NewRecorder()
	sh := newStaticHandler(logTesting)
	handler := http.HandlerFunc(sh.healthLive)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

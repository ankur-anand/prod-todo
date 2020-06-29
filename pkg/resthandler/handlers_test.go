// +build unit_tests all_tests

package resthandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankur-anand/prod-todo/pkg/logger"
)

func TestMuxHandler_ServeHTTP(t *testing.T) {
	t.Parallel()
	l, _ := logger.NewTesting(nil)
	defer l.Sync()
	mux := NewMuxHandler(l)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected error code %d for missing content type got %d", http.StatusBadRequest, rr.Code)
	}

	req.Header.Add("Content-type", "application/json")

	rr2 := httptest.NewRecorder()
	mux.ServeHTTP(rr2, req)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected error code %d for missing content type got %d", http.StatusOK, rr2.Code)
	}

	if rr2.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("mime sniff prevention header in response missing")
	}

	if rr2.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("content-type header in response is missing")
	}
}

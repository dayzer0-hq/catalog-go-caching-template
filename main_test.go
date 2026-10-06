package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The tests shrink the dataset and zero the simulated round-trip so they stay
// fast; the slow defaults exist for the running server, not for go test.
func testStore(t *testing.T) *Store {
	t.Setenv("SEED_ROWS", "100")
	t.Setenv("DB_DELAY_MS", "0")
	return NewStore()
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux(testStore(t)).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health = %d, want 200", rec.Code)
	}
}

func TestGetProduct(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux(testStore(t)).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/products/1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /products/1 = %d, want 200", rec.Code)
	}
	var p Product
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("response is not a product: %v", err)
	}
	if p.ID != 1 {
		t.Fatalf("got product %d, want 1", p.ID)
	}
}

func TestUnknownProductIs404(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux(testStore(t)).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/products/99999", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET /products/99999 = %d, want 404", rec.Code)
	}
}

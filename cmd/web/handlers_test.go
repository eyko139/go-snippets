package main

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "github.com/eyko139/go-snippets/internal/assert"
    "io"
)

func TestHealth(t *testing.T) {
    rr := httptest.NewRecorder()

    r, err := http.NewRequest(http.MethodGet, "/health", nil)

    if err != nil {
        t.Fatalf("Failed to execute request %s", err)
    }

    healthHandler := health()
    healthHandler(rr, r)
    
    rs := rr.Result()

    assert.AssertEqual(t, rs.StatusCode, http.StatusOK)

    defer rs.Body.Close()
    body, err := io.ReadAll(rs.Body)

    if err != nil {
        t.Fatalf("failed to parse body, err: %s", err)
    }

    assert.AssertEqual(t, string(body),"OK")

}

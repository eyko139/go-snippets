package main

import (
    "net/http/httptest"
    "net/http"
    "github.com/eyko139/go-snippets/internal/assert"
    "io"
    "testing"
)

func TestSecureHeaders(t *testing.T) {
    rr := httptest.NewRecorder()

    r, err := http.NewRequest(http.MethodGet, "/", nil)

    if err != nil {
        t.Fatalf("Failed to create Request, err: %s", err)
    }

    // Create a mock http.Handler which can be passed to the middleware

    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    }) 

    secureHeaders(next).ServeHTTP(rr, r)

    // check the result of the responseRecorder

    rs := rr.Result()

    expected := "nosniff"
    actual := rs.Header.Get("X-Content-Type-Options")
    assert.AssertEqual(t, actual, expected)

    expected = "deny"
    actual = rs.Header.Get("X-Frame-Options")
    assert.AssertEqual(t, actual, expected)

    expected = "0"
    actual = rs.Header.Get("X-XSS-Protection")
    assert.AssertEqual(t, actual, expected)


    assert.AssertEqual(t, rs.StatusCode, http.StatusOK)

    defer rs.Body.Close()

    body, err := io.ReadAll(rs.Body)

    if err != nil {
        t.Fatalf("Could not parse body, err: %s", err)
    }


    assert.AssertEqual(t, string(body), "OK")

}

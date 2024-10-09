package main

import (
	"github.com/eyko139/go-snippets/internal/assert"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {

	env := NewEnv()

	app, err := newTestApplication(env)

	ts := newTestServer(app.Routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/health")

	assert.AssertEqual(t, code, http.StatusOK)

	if err != nil {
		t.Fatal(err)
	}

	assert.AssertEqual(t, body, "OK")

}

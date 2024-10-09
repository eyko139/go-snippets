package main

import (
	"net/http"
	"testing"

	"github.com/eyko139/go-snippets/internal/assert"
	"github.com/eyko139/go-snippets/internal/session"
)

func TestHealth(t *testing.T) {

	env := NewEnv()

	app, err := newTestApplication(env)
	if err != nil {
		t.Fatal(err)
	}

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/health")

	assert.AssertEqual(t, code, http.StatusOK)
	assert.AssertEqual(t, body, "OK")
    session.DestroyProvider("memory")
}

func TestSnippetView(t *testing.T) {

	env := NewEnv()

	app, err := newTestApplication(env)

	if err != nil {
		t.Fatal(err)
	}

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "getFirst",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "mockContent",
		},
	}

    for _, test := range(tests) {
        t.Run(test.name, func(t *testing.T) {
            code, _, body := ts.get(t, test.urlPath)
            assert.AssertEqual(t, code, http.StatusOK)

            if test.wantBody != "" {
                assert.StringContains(t, body, test.wantBody)
            }
        })
    }
    session.DestroyProvider("memory")
}

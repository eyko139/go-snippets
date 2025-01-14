package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eyko139/go-snippets/cmd/util"
	"github.com/eyko139/go-snippets/internal/models"
	"github.com/eyko139/go-snippets/internal/models/mocks"
	"github.com/eyko139/go-snippets/internal/session"
	"github.com/eyko139/go-snippets/internal/session/providers"
)

func newTestApplication(env *Env) (*Config, error) {

	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	providers.InitMemorySession()
	globalSessions, err := session.NewManager("memory", "gosessionid", 360)

	if err != nil {
		errLog.Printf("Could not initialize session manager")
	}

	go globalSessions.GC()

	tc, err := models.NewTemplateCache()

	broker := models.NewBroker(env.BrokerConnection, infoLog, errLog)

	if err != nil {
		return nil, err
	}

	return &Config{
		ErrorLog:       errLog,
		InfoLog:        infoLog,
		Hlp:            util.NewHelper(tc, errLog, infoLog),
		GlobalSessions: globalSessions,
		Broker:         broker,
		Snippets:       &mocks.SnippetModel{},
		UserModel:      &mocks.UserModel{},
	}, nil
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, handler http.Handler) *testServer {
	ts := httptest.NewServer(handler)
	jar, err := cookiejar.New(nil)

	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar
	// Disable redirect-following for the test server client by setting a custom
	// CheckRedirect function. This function will be called whenever a 3xx
	// response is received by the client, and by always returning a
	// http.ErrUseLastResponse error it forces the client to immediately return
	// the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, path string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + path)

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

package mirait_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/catatsuy/bento/config"
	"github.com/catatsuy/bento/mirait"
)

func TestNewClient_parsesURL(t *testing.T) {
	defer mirait.SetTargetURL("https://example.com/foo/bar")()
	s, err := mirait.NewSession(config.Config{})
	if err != nil {
		t.Fatal(err)
	}

	expected := &url.URL{
		Scheme: "https",
		Host:   "example.com",
	}
	if !reflect.DeepEqual(s.URL, expected) {
		t.Fatalf("expected %q to equal %q", s.URL, expected)
	}
}

func TestGetToken_Success(t *testing.T) {
	muxAPI := http.NewServeMux()
	testAPIServer := httptest.NewServer(muxAPI)
	defer testAPIServer.Close()

	expectedUserAgent := "tester user agent"

	muxAPI.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if userAgent != expectedUserAgent {
			t.Fatalf("User-Agent: got %s, want %s", userAgent, expectedUserAgent)
		}

		if r.URL.Path != "/trial/" {
			t.Fatalf("got %s, want \"/trial/\"", r.URL.Path)
		}

		http.ServeFile(w, r, "testdata/get_token_ok.html")
	})

	defer mirait.SetTargetURL(testAPIServer.URL)()

	s, err := mirait.NewSession(config.Config{
		UserAgent: expectedUserAgent,
	})

	if err != nil {
		t.Fatal(err)
	}

	token, err := s.GetToken()
	if err != nil {
		t.Fatal(err)
	}

	expectedToken := "tokentokentoken"
	if token != expectedToken {
		t.Errorf("got %q, want %q", token, expectedToken)
	}
}

package mirait_test

import (
	"io/ioutil"
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

		if r.URL.Path != "/trial" {
			t.Fatalf("got %s, want \"/trial\"", r.URL.Path)
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

func TestPostTranslate_Success(t *testing.T) {
	muxAPI := http.NewServeMux()
	testAPIServer := httptest.NewServer(muxAPI)
	defer testAPIServer.Close()

	expectedToken := "tokentokentoken"
	expectedInput := "test"
	expectedIsJP := false

	muxAPI.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		expectedType := "application/x-www-form-urlencoded"
		if contentType != expectedType {
			t.Errorf("Content-Type expected %s, but %s", expectedType, contentType)
		}

		ck := r.Header.Get("Cookie")
		expectedCookie := "test=test_value"
		if ck != expectedCookie {
			t.Errorf("Cookie: got %s, want %s", ck, expectedCookie)
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()

		actualV, err := url.ParseQuery(string(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}

		expectedV := url.Values{}
		expectedV.Set("input", expectedInput)
		expectedV.Set("tran", expectedToken)

		if expectedIsJP {
			expectedV.Set("source", "ja")
			expectedV.Set("target", "en")
		} else {
			expectedV.Set("source", "en")
			expectedV.Set("target", "ja")
		}

		if !reflect.DeepEqual(actualV, expectedV) {
			t.Errorf("expected %q to equal %q", actualV, expectedV)
		}

		http.ServeFile(w, r, "testdata/post_translate_ok.json")
	})

	defer mirait.SetTargetURL(testAPIServer.URL)()

	s, err := mirait.NewSession(config.Config{})
	if err != nil {
		t.Fatal(err)
	}

	s.SetCacheCookie([]config.Cookie{
		config.Cookie{
			Name:  "test",
			Value: "test_value",
		},
	})
	s.SetToken(expectedToken)

	expectedIsJP = false
	output, err := s.PostTranslate(expectedInput, expectedIsJP)
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "こんにちは。"
	if output != expectedOutput {
		t.Errorf("got %s, want %s", output, expectedOutput)
	}

	expectedIsJP = true
	output, err = s.PostTranslate(expectedInput, expectedIsJP)
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput = "こんにちは。"
	if output != expectedOutput {
		t.Errorf("got %s, want %s", output, expectedOutput)
	}
}

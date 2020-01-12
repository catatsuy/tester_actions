package mirait

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultAPITimeout = 10
)

type Session struct {
	URL        *url.URL
	HTTPClient *http.Client

	Token string
}

func NewSession() (*Session, error) {
	jar, _ := cookiejar.New(&cookiejar.Options{})

	session := &Session{
		URL: &url.URL{
			Scheme: "https",
			Host:   "miraitranslate.com",
		},
		HTTPClient: &http.Client{
			Jar:     jar,
			Timeout: time.Duration(DefaultAPITimeout) * time.Second,
		},
	}

	return session, nil
}

func (s *Session) SetToken(ctx context.Context) error {
	u := s.URL
	u.Path = "/trial/"

	res, err := s.HTTPClient.Get(u.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse as HTML: %w", err)
	}

	token := ""
	doc.Find("input#tranInput").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		token, _ = s.Attr("value")
		return false
	})

	if token == "" {
		return errors.New("empty token")
	}

	s.Token = token

	return nil
}

func (s *Session) PostTranslate(ctx context.Context, input string) (output string, err error) {
	u := s.URL
	u.Path = "/trial/translate.php"

	q := url.Values{}
	q.Set("input", input)
	q.Set("tran", s.Token)
	q.Set("source", "en")
	q.Set("target", "ja")

	res, err := s.HTTPClient.PostForm(u.String(), q)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(bb), nil
}

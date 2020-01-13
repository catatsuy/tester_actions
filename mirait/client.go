package mirait

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/catatsuy/bento/config"
)

const (
	DefaultAPITimeout = 10

	userAgent = "waiwai client"
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

func (s *Session) SetCacheCookie(ccs []config.Cookie) {
	cookies := make([]*http.Cookie, 0, len(ccs))
	for _, cc := range ccs {
		cookies = append(cookies, &http.Cookie{
			Name:  cc.Name,
			Value: cc.Value,
		})
	}
	s.HTTPClient.Jar.SetCookies(s.URL, cookies)
}

func (s *Session) GetToken() (string, error) {
	u := s.URL
	u.Path = "/trial/"

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse as HTML: %w", err)
	}

	token := ""
	doc.Find("input#tranInput").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		token, _ = s.Attr("value")
		return false
	})

	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}

func (s *Session) SetToken(token string) {
	s.Token = token
}

type outputRes struct {
	Output string `json:"output"`
}

// {"status":"success","outputs":[{"output":"こんにちは。"}]}
type postTranslateRes struct {
	Status  string      `json:"status"`
	Outputs []outputRes `json:"outputs"`
}

func (s *Session) PostTranslate(input string, isJP bool) (output string, err error) {
	u := s.URL
	u.Path = "/trial/translate.php"

	q := url.Values{}
	q.Set("input", input)
	q.Set("tran", s.Token)

	if isJP {
		q.Set("source", "ja")
		q.Set("target", "en")
	} else {
		q.Set("source", "en")
		q.Set("target", "ja")
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(q.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	ptr := &postTranslateRes{}
	err = json.NewDecoder(res.Body).Decode(ptr)
	if err != nil {
		return "", fmt.Errorf("failed to encode json: %w", err)
	}

	if len(ptr.Outputs) == 0 {
		return "", fmt.Errorf("empty response")
	}

	if ptr.Status != "success" {
		return "", fmt.Errorf("no success response: %s", ptr.Status)
	}

	return ptr.Outputs[0].Output, nil
}

func (s *Session) DumpCookies() []config.Cookie {
	cookies := s.HTTPClient.Jar.Cookies(s.URL)
	ccs := make([]config.Cookie, 0, len(cookies))

	for _, c := range cookies {
		cc := config.Cookie{
			Name:  c.Name,
			Value: c.Value,
		}
		ccs = append(ccs, cc)
	}

	return ccs
}

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultAPITimeout = 10

	userAgent = "bento client"
)

func main() {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	c := &http.Client{
		Jar:     jar,
		Timeout: time.Duration(DefaultAPITimeout) * time.Second,
	}

	res, err := c.Get("https://miraitranslate.com/trial/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	tran := ""
	// Find the review items
	doc.Find("input#tranInput").Each(func(_ int, s *goquery.Selection) {
		// For each item found, get the band and title
		tran, _ = s.Attr("value")
	})

	q := url.Values{}
	q.Set("input", "Hello world")
	q.Set("tran", tran)
	q.Set("source", "en")
	q.Set("target", "ja")

	body := strings.NewReader(q.Encode())

	req, err := http.NewRequest("POST", "https://miraitranslate.com/trial/translate.php", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Origin", "https://miraitranslate.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(bb))
}

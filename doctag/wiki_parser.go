package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type WikiSearchAPI struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Continue string `json:"continue"`
		Sroffset int    `json:"sroffset"`
	} `json:"continue"`
	Query struct {
		Search []struct {
			Ns        int    `json:"ns"`
			Size      int    `json:"size"`
			Snippet   string `json:"snippet"`
			Timestamp string `json:"timestamp"`
			Title     string `json:"title"`
			Wordcount int    `json:"wordcount"`
		} `json:"search"`
		Searchinfo struct {
			Totalhits int `json:"totalhits"`
		} `json:"searchinfo"`
	} `json:"query"`
}

const apiURL = "https://en.wikipedia.org/w/api.php?action=query&list=search&utf8=&format=json"
const wikiURL = "https://en.wikipedia.org/wiki/"

func getWikipediaReader(keywords []string) (io.Reader, error) {
	u, _ := url.Parse(apiURL)
	qval := u.Query()
	qval.Add("srsearch", strings.Join(keywords, " "))
	u.RawQuery = qval.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	w := WikiSearchAPI{}
	if err := json.Unmarshal(bs, &w); err != nil {
		return nil, err
	}

	if len(w.Query.Search) < 1 {
		return nil, errors.New("No results found.")
	}
	url := wikiURL + w.Query.Search[0].Title
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	str := getTextToken(bytes.NewReader(bs))
	return bytes.NewReader([]byte(str)), nil
}

func getTextToken(rd io.Reader) string {
	z := html.NewTokenizer(rd)
	ret := ""
	isScript := false
	isStyle := false
	done := false
	for !done {
		tt := z.Next()
		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			d := z.Token().Data
			if d == "script" {
				isScript = true
			}
			if d == "style" {
				isStyle = true
			}
		case html.TextToken:
			if !isScript && !isStyle {
				t := z.Token()
				ret += t.Data
			}
		case html.ErrorToken:
			break
		case html.EndTagToken:
			t := z.Token().Data
			if t == "html" {
				done = true
				break
			}
			if t == "script" {
				isScript = false
			}
			if t == "style" {
				isStyle = false
			}
		}
	}
	return ret
}

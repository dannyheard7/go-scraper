package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

func Fetch(url *url.URL) (string, []*url.URL, error) {
	siteBase, err := url.Parse(url.Scheme + "://" + url.Host)
	if err != nil {
		return "", nil, err
	}

	siteRules, err := get_site_rules(siteBase)
	if err != nil {
		return "", nil, err
	}

	// TODO: move this into the 'main' for loop
	for _, disallowedPage := range siteRules.disallowedPages {
		if disallowedPage.MatchString(url.String()) {
			return "", nil, errors.New("this page is not allowed to be indexed by the sites robot.txt")
		}
	}

	log.Debug().Msg("Can Index " + url.String())

	resp, err := http.Get(url.String())
	if err != nil {
		return "", nil, err
	}

	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("expected 200 status code, received %d", resp.StatusCode)
	}

	tkn := html.NewTokenizer(resp.Body)

	var body string
	var shouldRead bool

	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return body, nil, nil
		case tt == html.StartTagToken:
			shouldRead = is_text_token(tkn.Token())
		case tt == html.TextToken:
			if shouldRead {
				body = body + tkn.Token().Data
			}
		case tt == html.EndTagToken:
			if shouldRead && is_text_token(tkn.Token()) {
				shouldRead = false
				body = body + "\n"
			}
		}
	}
}

func is_text_token(token html.Token) bool {
	println(token.String())
	// TODO: this doesn't seem to work very well for nested data. e.g. <B> in <p> (WRITE tests)
	return token.Data == "p" || token.Data == "b" || strings.HasPrefix(token.Data, "h")
}

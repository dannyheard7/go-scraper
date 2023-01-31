package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

const UserAgentPrefix = "user-agent"
const AllowedPrefix = "allow"
const DisallowedPrefix = "disallow"
const Wildcard = "*"

type SiteScraperRules struct {
	allowedPages    []*regexp.Regexp
	disallowedPages []*regexp.Regexp
}

func get_site_rules(siteBase *url.URL) (*SiteScraperRules, error) {
	scanner, err := get_robots_txt(siteBase)
	if err != nil {
		return nil, err
	}

	return parse_robots_txt(siteBase, scanner)
}

func get_robots_txt(siteBase *url.URL) (*bufio.Scanner, error) {
	if len(siteBase.Path) > 0 {
		return nil, errors.New("expected a url with no path")
	}

	joinedPath, err := url.JoinPath(siteBase.String(), "robots.txt")
	if err != nil {
		return nil, err
	}

	response, err := http.Get(joinedPath)
	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(response.Body), nil
}

func parse_robots_txt(siteBase *url.URL, scanner *bufio.Scanner) (*SiteScraperRules, error) {
	robotsTxtResult := SiteScraperRules{allowedPages: make([]*regexp.Regexp, 0), disallowedPages: make([]*regexp.Regexp, 0)}

	relevant := false
	for scanner.Scan() {
		lineText := strings.TrimSpace(scanner.Text())
		if len(lineText) == 0 {
			log.Debug().Msg("Skipping empty line")
			continue
		}

		splits := strings.SplitN(lineText, ":", 2)

		if len(splits) != 2 {
			log.Debug().Msg(fmt.Sprintf("Failed to read line from robots.txt %s", lineText))
			continue
		}

		prefix, value := strings.TrimSpace(strings.ToLower(splits[0])), strings.TrimSpace(splits[1])

		if prefix == UserAgentPrefix && value == Wildcard {
			relevant = true
			continue
		}

		if prefix == UserAgentPrefix {
			relevant = false
			continue
		}

		if !relevant {
			continue
		}

		fullPath, err := url.JoinPath(siteBase.String(), value)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Cannot create full path from %s and %s, Error %s", fullPath, value, err))
			continue
		}

		regex, err := build_regex_from_robots_txt_value(fullPath)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Failed to convert value %s into regex, Error %s", fullPath, err))
			continue
		}

		if prefix == AllowedPrefix {
			robotsTxtResult.allowedPages = append(robotsTxtResult.allowedPages, regex)
		}

		if prefix == DisallowedPrefix {
			robotsTxtResult.disallowedPages = append(robotsTxtResult.disallowedPages, regex)
		}
	}

	return &robotsTxtResult, nil
}

func build_regex_from_robots_txt_value(value string) (*regexp.Regexp, error) {
	escaped := strings.ReplaceAll(value, "/", "\\/")
	escaped = strings.ReplaceAll(escaped, "?", "\\?")
	escaped = strings.ReplaceAll(escaped, ".", "\\.")

	return regexp.Compile(escaped)
}

package main

import (
	"net/url"
	"web-scraper/crawler"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	baseUrl, err := url.Parse("https://www.bbc.co.uk/news/uk-england-merseyside-64454778")

	if err != nil {
		panic(err)
	}

	body, urls, err := crawler.Fetch(baseUrl)

	if err != nil {
		log.Error().Msg(err.Error())
	}

	println(body, urls)
}

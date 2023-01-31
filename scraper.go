package main

import (
	"net/url"
	"web-scraper/crawler"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	baseUrl, err := url.Parse("https://blog.logrocket.com/5-structured-logging-packages-for-go")

	if err != nil {
		panic(err)
	}

	body, urls, err := crawler.Fetch(baseUrl)

	if err != nil {
		log.Error().Msg(err.Error())
	}

	println(body, urls)
}

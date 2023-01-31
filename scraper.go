package main

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	baseUrl, err := url.Parse("https://google.com")

	if err != nil {
		panic(err)
	}

	resp, err := get_site_rules(baseUrl)
	fmt.Println(resp.allowedPages)
	// fmt.Println(resp.disallowedPages)
}

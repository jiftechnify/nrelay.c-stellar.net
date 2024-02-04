package main

import (
	"log"
	"os"

	"github.com/jiftechnify/strfrui"
	"github.com/jiftechnify/strfrui/sifters"

	"c-stellar-relay-evsifter/api"
)

func main() {
	apiBaseURL := os.Getenv("FOLLOW_CHECK_API_BASE_URL")
	if apiBaseURL == "" {
		log.Fatal("FOLLOW_CHECK_API_BASE_URL is not set")
	}
	apiCli := api.NewClient(apiBaseURL)

	sifter := sifters.AuthorMatcher(apiCli.IsFollower, sifters.Allow)
	strfrui.New(sifter).Run()
}

package main

import (
	"log"
	"os"

	"github.com/jiftechnify/strfrui"
	"github.com/jiftechnify/strfrui/sifters"
	"github.com/jiftechnify/strfrui/sifters/ratelimit"

	"c-stellar-relay-evsifter/api"
)

// 500 events/hr + burst of 50 events. Ephemeral events are unlimited.
var basicRatelimit = ratelimit.
	ByUser(ratelimit.Quota{
		MaxRate:  ratelimit.PerHour(500),
		MaxBurst: 50,
	}, ratelimit.PubKey).
	Exclude(func(i *strfrui.Input) bool { return sifters.KindsAllEphemeral(i.Event.Kind) })

func main() {
	apiBaseURL := os.Getenv("FOLLOW_CHECK_API_BASE_URL")
	if apiBaseURL == "" {
		log.Fatal("FOLLOW_CHECK_API_BASE_URL is not set")
	}
	apiCli := api.NewClient(apiBaseURL)

	sifter := sifters.Pipeline(
		sifters.AuthorMatcher(apiCli.IsFollower, sifters.Allow),
		basicRatelimit,
	)

	strfrui.New(sifter).Run()
}

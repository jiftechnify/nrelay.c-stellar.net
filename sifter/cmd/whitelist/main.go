package main

import (
	"log"
	"os"

	"github.com/jiftechnify/strfrui"
	"github.com/jiftechnify/strfrui/sifters"
	"github.com/jiftechnify/strfrui/sifters/ratelimit"

	"c-stellar-relay-evsifter/api"
)

var limitReactionsPerMin = ratelimit.ByUserAndKind([]ratelimit.KindQuota{
	ratelimit.QuotaForKinds([]int{7}, ratelimit.Quota{MaxRate: ratelimit.PerMin(1), MaxBurst: 1}),
}, ratelimit.PubKey)

func main() {
	apiBaseURL := os.Getenv("FOLLOW_CHECK_API_BASE_URL")
	if apiBaseURL == "" {
		log.Fatal("FOLLOW_CHECK_API_BASE_URL is not set")
	}
	apiCli := api.NewClient(apiBaseURL)

	sifter := sifters.Pipeline(
		sifters.AuthorMatcher(apiCli.IsFollower, sifters.Allow),
		limitReactionsPerMin,
	)

	strfrui.New(sifter).Run()
}

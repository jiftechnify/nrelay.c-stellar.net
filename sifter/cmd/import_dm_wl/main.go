package main

import (
	"c-stellar-relay-evsifter/api"
	"log"
	"os"

	"github.com/jiftechnify/strfrui"
	"github.com/jiftechnify/strfrui/sifters"
	"github.com/nbd-wtf/go-nostr"
)

func allowIfRecipientIsMyFollower(apiCli *api.Client) strfrui.Sifter {
	return sifters.TagsMatcher(func(tags nostr.Tags) (bool, error) {
		for _, t := range tags.GetAll([]string{"p", ""}) {
			isFollower, err := apiCli.IsFollower(t.Value())
			if err != nil {
				return false, err
			}
			if isFollower { // if any of the recipients is my follower, allow
				return true, nil
			}
		}
		return false, nil
	}, sifters.Allow)
}

func main() {
	apiBaseURL := os.Getenv("FOLLOW_CHECK_API_BASE_URL")
	if apiBaseURL == "" {
		log.Fatal("FOLLOW_CHECK_API_BASE_URL is not set")
	}
	apiCli := api.NewClient(apiBaseURL)

	sifter := sifters.Pipeline(
		// allow only DMs
		sifters.KindList([]int{4}, sifters.Allow).RejectWithMsg("blocked: not DM"),
		sifters.OneOf(
			// accept DMs from white-listed pubkeys
			sifters.AuthorMatcher(apiCli.IsFollower, sifters.Allow),
			// accept DMs to white-listed pubkeys
			allowIfRecipientIsMyFollower(apiCli),
		).RejectWithMsg("blocked: sender nor recipient is in the whitelist"),
	)

	strfrui.New(sifter).Run()
}

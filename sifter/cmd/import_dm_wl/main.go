package main

import (
	"c-stellar-relay-evsifter/api"

	"github.com/jiftechnify/strfrui"
	"github.com/jiftechnify/strfrui/sifters"
	"github.com/nbd-wtf/go-nostr"
)

var allowIfRecipientIsMyFollower = sifters.TagsMatcher(func(tags nostr.Tags) (bool, error) {
	for _, t := range tags.GetAll([]string{"p"}) {
		isFollower, err := api.IsFollower(t.Value())
		if err != nil {
			return false, err
		}
		if isFollower { // if any of the recipients is my follower, allow
			return true, nil
		}
	}
	return false, nil
}, sifters.Allow)

func main() {
	sifter := sifters.Pipeline(
		// allow only DMs
		sifters.KindList([]int{4}, sifters.Allow).RejectWithMsg("blocked: not DM"),
		sifters.OneOf(
			// accept DMs from white-listed pubkeys
			sifters.AuthorMatcher(api.IsFollower, sifters.Allow),
			// accept DMs to white-listed pubkeys
			allowIfRecipientIsMyFollower,
		).RejectWithMsg("blocked: sender nor recipient is in the whitelist"),
	)

	strfrui.New(sifter).Run()
}

package main

import (
	"github.com/jiftechnify/strfrui"
	"github.com/jiftechnify/strfrui/sifters"

	"c-stellar-relay-evsifter/api"
)

var followersOnlySifter = sifters.AuthorMatcher(api.IsFollower, sifters.Allow)

func main() {
	strfrui.New(followersOnlySifter).Run()
}

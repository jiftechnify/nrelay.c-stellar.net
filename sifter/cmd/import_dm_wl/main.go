package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/jiftechnify/strfrui"
	"github.com/jiftechnify/strfrui/sifters"
	"github.com/nbd-wtf/go-nostr"
)

func main() {
	resDir := os.Getenv("RESOURCE_DIR")
	if resDir == "" {
		log.Fatal("RESOURCE_DIR is not set")
	}

	wlSifter, err := readWhiteList(filepath.Join(resDir, "whitelist.txt"))
	if err != nil {
		log.Fatal(err)
	}

	strfrui.New(wlSifter).Run()
}

func readWhiteList(path string) (strfrui.Sifter, error) {
	wlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer wlFile.Close()

	wl := make([]string, 0)
	scanner := bufio.NewScanner(wlFile)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Print(err)
			continue
		}
		wl = append(wl, scanner.Text())
	}

	s := sifters.Pipeline(
		// allow only DMs
		sifters.KindList([]int{4}, sifters.Allow).RejectWithMsg("blocked: not DM"),
		sifters.OneOf(
			// accept DMs from white-listed pubkeys
			sifters.AuthorList(wl, sifters.Allow),
			// accept DMs to white-listed pubkeys
			sifters.TagsMatcher(func(t nostr.Tags) (bool, error) {
				return t.ContainsAny("p", wl), nil
			}, sifters.Allow),
		).RejectWithMsg("blocked: sender nor recipient is in the whitelist"),
	)
	return s, nil
}

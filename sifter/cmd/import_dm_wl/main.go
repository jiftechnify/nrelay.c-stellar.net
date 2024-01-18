package main

import (
	"bufio"
  "errors"
	"log"
	"os"

	evsifter "github.com/jiftechnify/strfry-evsifter"
)

func getDestination(input *evsifter.Input) (string, error) {
  for _, tag := range input.Event.Tags {
    if len(tag) <= 1 || tag[0] != "p" || len(tag[1]) != 64  {
      continue
    }
    return tag[1], nil
  }
  return "", errors.New("DM destination not found")
}

type ImportDMWithWhiteListSifer map[string]struct{}

func (wl ImportDMWithWhiteListSifer) Sift(input *evsifter.Input) (*evsifter.Result, error) {
  if input.Event.Kind != 4 {
    return input.Reject("blocked: not DM")
  }

  // accept DMs from white-listed pubkeys
  if _, ok := wl[input.Event.PubKey]; ok {
    return input.Accept()
  }

  // accept DMs to white-listed pubkeys
  dest, err := getDestination(input)
  if err != nil {
    return input.Reject("blocked: unknown DM destination")
  }
	if _, ok := wl[dest]; ok {
		return input.Accept()
	}

	return input.Reject("blocked: the destination pubkey is not in the white-list")
}

func readWhiteList(path string) (ImportDMWithWhiteListSifer, error) {
	wlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer wlFile.Close()

	wl := make(ImportDMWithWhiteListSifer, 0)
	scanner := bufio.NewScanner(wlFile)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Print(err)
			continue
		}
		wl[scanner.Text()] = struct{}{}
	}
	return wl, nil
}

func main() {
	wl, err := readWhiteList("./resource/whitelist.txt")
	if err != nil {
		log.Fatal(err)
	}

	var s evsifter.Runner
	s.SiftWith(wl)
	s.Run()
}

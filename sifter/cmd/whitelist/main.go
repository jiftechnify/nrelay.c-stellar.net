package main

import (
	"bufio"
	"log"
	"os"

	evsifter "github.com/jiftechnify/strfry-evsifter"
)

type WhiteListSifer map[string]struct{}

func (w WhiteListSifer) Sift(input *evsifter.Input) (*evsifter.Result, error) {
	if _, ok := w[input.Event.PubKey]; ok {
		return input.Accept()
	}
	return input.Reject("blocked: you can't write events")
}

func readWhiteList(path string) (WhiteListSifer, error) {
	wlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer wlFile.Close()

	wl := make(WhiteListSifer, 0)
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

package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/jiftechnify/strfrui"
	"github.com/jiftechnify/strfrui/sifters"
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

	return sifters.AuthorList(wl, sifters.Allow), nil
}

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	cache   map[string]time.Time
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		cache:   make(map[string]time.Time),
	}
}

type isFollowerResp struct {
	IsFollower bool `json:"isFollower"`
}

const cacheTTL = 1 * time.Hour

func (c *Client) IsFollower(author string) (bool, error) {
	if t, ok := c.cache[author]; ok && time.Since(t) < cacheTTL {
		log.Printf("is_follower cache hit (pubkey: %s)", author)
		return true, nil
	}

	apiURL := fmt.Sprintf("%s/is_follower?pubkey=%s", c.baseURL, author)
	resp, err := http.Get(apiURL)
	if err != nil {
		return false, fmt.Errorf("API server error: %w", err)
	}
	defer resp.Body.Close()

	var r isFollowerResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, fmt.Errorf("failed to decode response from API server: %w", err)
	}

	// update cache
	if r.IsFollower {
		c.cache[author] = time.Now()
	} else {
		delete(c.cache, author)
	}
	return r.IsFollower, nil
}

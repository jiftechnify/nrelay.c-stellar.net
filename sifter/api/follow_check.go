package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: strings.TrimSuffix(baseURL, "/")}
}

type isFollowerResp struct {
	IsFollower bool `json:"is_follower"`
}

func (c *Client) IsFollower(author string) (bool, error) {
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
	return r.IsFollower, nil
}

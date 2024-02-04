package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type isFollowerResp struct {
	IsFollower bool `json:"is_follower"`
}

func IsFollower(author string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("http://whitelisted_follow_check_api:8080/is_follower?pubkey=$%s", author))
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

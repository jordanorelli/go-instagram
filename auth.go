package instagram

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

func (c *Client) RequestToken(code string) (*AuthResponse, error) {
	path := "https://api.instagram.com/oauth/access_token"

	raw, err := http.PostForm(path, url.Values{
		"client_id":     []string{c.clientId},
		"client_secret": []string{c.clientSecret},
		"grant_type":    []string{"authorization_code"},
		"redirect_uri":  []string{c.redirectUri},
		"code":          []string{code},
	})
	if err != nil {
		return nil, err
	}
	defer raw.Body.Close()

	var response AuthResponse
	if err := json.NewDecoder(raw.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

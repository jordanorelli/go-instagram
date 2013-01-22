package instagram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const baseUrl = "https://api.instagram.com/v1"

type Client struct {
	*http.Client
	clientId     string
	clientSecret string
	redirectUri  string
}

func NewClient(clientId, clientSecret, redirectUri string) *Client {
	return &Client{new(http.Client), clientId, clientSecret, redirectUri}
}

type Envelope interface {
	GetMeta() *Meta
}

type APIError struct {
	Code         int    `json:"code"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf(`instagram: %d %v: %v`, e.Code, e.ErrorType, e.ErrorMessage)
}

func unmarshalResponse(res *http.Response, v Envelope) error {
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		var m struct {
			Meta Meta `json:"meta"`
		}
		if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
			return fmt.Errorf(`instagram: unable to unmarshal API response: %v`, err)
		}
		return &APIError{
			Code:         m.Meta.Code,
			ErrorType:    m.Meta.ErrorType,
			ErrorMessage: m.Meta.ErrorMessage,
		}
	}
	return json.NewDecoder(res.Body).Decode(v)
}

// makes an api request at the given path, using the provided access token.  on
// success, the json body is unmarshaled into the provided envelope.
func (c *Client) get(path string, token string, v Envelope) error {
	var url string
	if token != "" {
		url = baseUrl + path + "?access_token=" + token
	} else {
		url = baseUrl + path + "?client_id=" + c.clientId
	}
	log.Println(url)
	res, err := c.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return unmarshalResponse(res, v)
}

type Meta struct {
	Code         int    `json:"code"`
	ErrorType    string `json:"error_type,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type Pagination struct {
	NextUrl    string `json:"next_url"`
	NextMaxId  string `json:"next_max_id,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

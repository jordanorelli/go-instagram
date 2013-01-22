package instagram

import (
	"net/http"
)

type User struct {
	Bio            string `json:"bio,omitempty"`
	FullName       string `json:"full_name"`
	Id             string `json:"id" bson:"id"`
	ProfilePicture string `json:"profile_picture"`
	Username       string `json:"username"`
	Website        string `json:"website,omitempty"`
}

type UserPage struct {
	Meta       *Meta       `json:"meta"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Data       []*User     `json:"data"`
}

// this is bullshit!  it's a bug in encoding/json.
func (p *UserPage) GetMeta() *Meta { return p.Meta }

func (p *UserPage) NextUrl() string { return p.Pagination.NextUrl }
func (p *UserPage) HasNext() bool   { return p.Pagination.NextUrl != "" }

func (page *UserPage) Next() (*UserPage, error) {
	c := &http.Client{}
	res, err := c.Get(page.Pagination.NextUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var users UserPage
	if err := unmarshalResponse(res, &users); err != nil {
		return nil, err
	}
	return &users, nil
}

func (page *UserPage) Ids() []string {
	ids := make([]string, len(page.Data))
	for i, user := range page.Data {
		ids[i] = user.Id
	}
	return ids
}

type FollowingPager struct {
	UserId string
	Token  string
	err    error
}

// given a user ID and an optional oauth token, retrieves the user's profile info from the Instagram API.
func (c *Client) User(userId string, token string) (*User, error) {
	var res userResponse
	if err := c.get("/users/"+userId, token, &res); err != nil {
		return nil, err
	}
	return res.User, nil
}

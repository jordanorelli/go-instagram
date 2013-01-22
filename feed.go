package instagram

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

type FeedPage struct {
	Meta       *Meta       `json:"meta"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Data       []Media     `json:"data"`

	client *Client `json:"-"`
}

// this is bullshit!
func (p *FeedPage) GetMeta() *Meta { return p.Meta }

func (p *FeedPage) NextUrl() string {
	if p.Pagination == nil {
		return ""
	}
	return p.Pagination.NextUrl
}

func (p *FeedPage) HasNext() bool {
	return p.Pagination != nil && p.Pagination.NextUrl != ""
}

type FeedParam struct {
	Token string
	MinId string
	Count int
	Reply chan *FeedPage
}

type FeedPager struct {
	Token string
	MinId string
	err   error
}

func (c *Client) getFeedPage(url string, response *FeedPage) error {
	*response = FeedPage{client: c}
	log.Println(url)
	res, err := c.Client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = unmarshalResponse(res, response)
	if err != nil {
		return err
	}
	return nil
}

func feedUrl(params map[string]string) string {
	v := url.Values{}
	for key, val := range params {
		v.Set(key, val)
	}
	return baseUrl + "/users/self/feed?" + v.Encode()
}

func (c *Client) Feed(token string, count int, response *FeedPage) error {
	return c.getFeedPage(feedUrl(map[string]string{
		"access_token": token,
		"count":        strconv.Itoa(count),
	}), response)
}

func (c *Client) FeedBefore(token string, count int, response *FeedPage, maxId string) error {
	return c.getFeedPage(feedUrl(map[string]string{
		"access_token": token,
		"count":        strconv.Itoa(count),
		"max_id":       maxId,
	}), response)
}

func (c *Client) FeedAfter(token string, count int, response *FeedPage, minId string) error {
	return c.getFeedPage(feedUrl(map[string]string{
		"access_token": token,
		"count":        strconv.Itoa(count),
		"min_id":       minId,
	}), response)
}

func (p *FeedPage) Next() error {
	if p.client == nil {
		return errors.New("instagram: nil client in FeedPage.Next()")
	}
	return p.client.getFeedPage(p.NextUrl(), p)
}

func (page *FeedPage) NextMinId() string {
	if len(page.Data) == 0 {
		return ""
	}
	return page.Data[0].Id
}

func mkUrl(token, minId, maxId string) string {
	v := url.Values{
		"access_token": []string{token},
	}
	if minId != "" {
		v["min_id"] = []string{minId}
	}
	if maxId != "" {
		v["max_id"] = []string{maxId}
	}
	return fmt.Sprintf("%s/users/self/feed?%s", baseUrl, v.Encode())
}

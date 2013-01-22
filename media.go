package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// used in both caption and comments
type ShortText struct {
	CreatedTimestamp *Timestamp `json:"created_time"`
	From             *User      `json:"from"`
	Id               string     `json:"id"`
	Text             string     `json:"text"`
}

type Comments struct {
	Count int          `json:"count"`
	Data  []*ShortText `json:"data"`
}

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type ImageSet struct {
	Low       *Image `json:"low_resolution"`
	Standard  *Image `json:"standard_resolution"`
	Thumbnail *Image `json:"thumbnail"`
}

type Likes struct {
	Data  []*User `json:"data"`
	Count int     `json:"count"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Id        int     `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
}

type Media struct {
	Id               string     `json:"id" bson:"_id"`
	Caption          *ShortText `json:"caption"`
	Comments         *Comments  `json:"comments"`
	CreatedTimestamp *Timestamp `json:"created_time"`
	Filter           NString    `json:"filter"`
	Images           *ImageSet  `json:"images"`
	Likes            Likes      `json:"likes"`
	Link             NString    `json:"link" bson:",omitempty"`
	Location         *Location  `json:"location"`
	Tags             []string   `json:"tags"`
	Type             string     `json:"type"`
	User             *User      `json:"user"`
}

//------------------------------------------------------------------------------
// this section of type definitions is clearly wrong.  Refactor the paging
// mechanics to use an interface.
//------------------------------------------------------------------------------

type MediaPage FeedPage

var ErrLastPage = errors.New("no more pages to return")

func (p *MediaPage) Next() (*MediaPage, error) {
	if p.Pagination == nil || p.Pagination.NextUrl == "" {
		return nil, ErrLastPage
	}
	if p.client == nil {
		return nil, errors.New("page has no instagram client")
	}
	return p.client.getMediaPage(p.Pagination.NextUrl)
}

func (p *MediaPage) GetMeta() *Meta { return p.Meta }

type MediaPager FeedPager

//------------------------------------------------------------------------------

type userResponse struct {
	Meta *Meta `json:"meta"`
	User *User `json:"data"`
}

func (res *userResponse) GetMeta() *Meta { return res.Meta }

func (c *Client) getMediaPage(url string) (*MediaPage, error) {
	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var page MediaPage
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, err
	}

	page.client = c
	return &page, nil
}

func (c *Client) GetUserMedia(userId, token string, count int) (*MediaPage, error) {
	url := fmt.Sprintf("%s/users/%s/media/recent?count=%d&access_token=%s", baseUrl, userId, count, token)
	return c.getMediaPage(url)
}

func (c *Client) GetUserMediaBefore(userId, token string, count int, before time.Time) (*MediaPage, error) {
	url := fmt.Sprintf("%s/users/%s/media/recent?count=%d&access_token=%s&max_timestamp=%d", baseUrl, userId, count, token, before.Unix())
	return c.getMediaPage(url)
}

func (c *Client) GetUserMediaBeforeId(userId, token string, count int, before string) (*MediaPage, error) {
	url := fmt.Sprintf("%s/users/%s/media/recent?count=%d&access_token=%s&max_id=%s", baseUrl, userId, count, token, before)
	return c.getMediaPage(url)
}

// given a userId of a target user whose media we would like to download, and
// an API token to use, continually gets more pages of media from the instagram
// API, putting it into a channel.
func (c *Client) StreamUserMedia(userId, token string, out chan *MediaPage, quit chan bool) {
	latestUrl := fmt.Sprintf("%s/users/%s/media/recent?count=100&access_token=%s", baseUrl, userId, token)
	url := latestUrl
	defer close(out)
	for {
		select {
		case <-quit:
			return
		default:
			page, err := c.getMediaPage(url)
			if err != nil {
				log.Println("ERROR paging through user's media history:", err.Error())
				return
			}
			out <- page
			if page.Pagination != nil && page.Pagination.NextUrl != "" {
				url = page.Pagination.NextUrl
				delay := time.Duration(1.5e9)
				time.Sleep(delay)
			} else {
				log.Println("Hit the end of the user's media history...")
				return
			}
		}
	}
}

func (p *FollowingPager) Error() string {
	if p.err == nil {
		return ""
	}
	return p.err.Error()
}

func (p *FollowingPager) Page(out chan *UserPage, quit chan bool) {
	c := new(http.Client)
	url := fmt.Sprintf("%s/users/%s/follows?access_token=%s", baseUrl, p.UserId, p.Token)

	defer close(out)
	for {
		select {
		case <-quit:
			return
		default:
			var page UserPage
			res, err := c.Get(url)
			if err != nil {
				log.Println("ERROR paging through user follows API:", err.Error())
				p.err = err
				return
			}
			defer res.Body.Close()

			err = unmarshalResponse(res, &page)
			if err != nil {
				p.err = err
				return
			}

			out <- &page
			if !page.HasNext() {
				return
			} else {
				url = page.Pagination.NextUrl
			}
		}
	}
}

func (m *Media) String() string {
	return fmt.Sprintf("[instagram_media Id: %s, Link: %s]", m.Id, m.Link)
}

func (u *User) String() string {
	return fmt.Sprintf("[instagram_user Id: %s, Username: %s]", u.Id, u.Username)
}

// GetMediaById returns a single media item identified by mediaId using the supplied token
func (c *Client) GetMediaById(mediaId, token string) (*MediaPage, error) {
	url := fmt.Sprintf("%s/media/%s?access_token=%s", baseUrl, mediaId, token)
	return c.getMediaPage(url)
}

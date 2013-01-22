package instagram

func (c *Client) GetTagItems(token, tagname string) ([]Media, error) {
	url := baseUrl + "/tags/" + tagname + "/media/recent"
	var page FeedPage
	err := c.get(url, token, &page)
	if err != nil {
		return nil, err
	}
	return page.Data, nil
}

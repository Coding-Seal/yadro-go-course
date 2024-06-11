package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var ErrNotAuthorized = errors.New("not authorized")

type Client struct {
	baseURL string
	c       *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		c: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Login(ctx context.Context, login, pswd string) (string, error) {
	body := bytes.NewBuffer(nil)

	err := json.NewEncoder(body).Encode(map[string]string{"login": login, "password": pswd})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/login", body)
	if err != nil {
		return "", err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrNotAuthorized
	}
	token := resp.Header.Get("Authorization")
	return token, nil
}

func (c *Client) SearchPics(ctx context.Context, query string) ([]string, error) {
	var pics []string

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/pics", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("search", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&pics)
	if err != nil {
		return nil, err
	}

	return pics, nil
}

func (c *Client) Update(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/update", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

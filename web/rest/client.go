package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrNotAuthorized = errors.New("rest api: not authorized")
	ErrInternal      = errors.New("rest api: internal error ")
)

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
		return "", errors.Join(err, ErrInternal)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/login", body)
	if err != nil {
		return "", errors.Join(err, ErrInternal)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return "", ErrNotAuthorized
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrNotAuthorized
	}

	token := resp.Header.Get("Authorization")

	return token, nil
}

func (c *Client) SearchPics(ctx context.Context, query string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/pics", nil)
	if err != nil {
		return nil, errors.Join(err, ErrInternal)
	}

	q := req.URL.Query()
	q.Add("search", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, errors.Join(err, ErrInternal)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w %s", ErrInternal, http.StatusText(resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)

	var pics []string

	err = decoder.Decode(&pics)
	if err != nil {
		return nil, errors.Join(err, ErrInternal)
	}

	return pics, nil
}

func (c *Client) Update(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/update", nil)
	if err != nil {
		return errors.Join(err, ErrInternal)
	}

	req.Header.Add("Authorization", token)

	resp, err := c.c.Do(req)
	if err != nil {
		return errors.Join(err, ErrInternal)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w %s", ErrInternal, http.StatusText(resp.StatusCode))
	}

	return nil
}

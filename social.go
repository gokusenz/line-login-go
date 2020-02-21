package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// A Client manages communication.
type Client struct {
	http *http.Client
	code string
}

// NewClient returns a new client. If a nil httpClient is
// Provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client, code string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		http: httpClient,
	}

	return c
}

// Do sends a success request to the website, a error is returned
func (c *Client) Do(ctx context.Context) (string, error) {
	return c.post(ctx)
}

func (c *Client) post(ctx context.Context) (string, error) {

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", c.code)
	data.Set("redirect_uri", "https://line-login-soical.herokuapp.com")
	data.Set("client_id", "1653859637")
	data.Set("client_secret", "350c50d8c1e4435726d64450f45142c3")

	req, err := http.NewRequest(http.MethodPost, "https://api.line.me/oauth2/v2.1/token", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return "", err
	}

	res, err := c.http.Do(req)
	if err != nil {
		log.Println(err)
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	s := string(string(body))

	return s, nil

}

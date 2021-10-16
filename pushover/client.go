package pushover

import (
	"context"
	"fmt"
	"github.com/hmoragrega/fastlane"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	user       string
	token      string
	baseURL    string
	sound      string
	httpClient *http.Client
}

type ClientOpt func(c *Client)

func WithHttpClient(httpClient *http.Client) ClientOpt {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithSound(sound string) ClientOpt {
	return func(c *Client) {
		c.sound = sound
	}
}

func New(user, token, baseURL string, opts ...ClientOpt) *Client {
	c := &Client{
		user:    user,
		token:   token,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
	for _, o := range opts {
		o(c)
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c
}

type pushRequest struct {
	Token    string `json:"token"`
	User     string `json:"user"`
	Message  string `json:"message"`
	HTML     string `json:"html,omitempty"`
	Sound    string `json:"sound,omitempty"`
	URL      string `json:"url,omitempty"`
	URLTitle string `json:"url_title,omitempty"`
}

func (c *Client) Push(ctx context.Context, message string, opts fastlane.PushOptions) error {
	v := make(url.Values)
	v.Set("user", c.user)
	v.Set("token", c.token)
	v.Set("message", message)
	if opts.HTML {
		v.Set("html", "1")
	}
	if opts.Sound != "" {
		v.Set("sound", opts.Sound)
	}
	if opts.URL != "" {
		v.Set("url", opts.URL)
	}
	if opts.URLTitle != "" {
		v.Set("url_title", opts.URLTitle)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/1/messages.json", strings.NewReader(v.Encode()))
	if err != nil {
		return fmt.Errorf("cannot create push request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot execute push request: %w", err)
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("push failed: req %q \n res %q", v.Encode(), string(b))
	}

	return nil
}

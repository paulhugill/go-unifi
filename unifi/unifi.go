package unifi //

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type NotFoundError struct{}

func (err *NotFoundError) Error() string {
	return "not found"
}

type APIError struct {
	RC      string
	Message string
}

func (err *APIError) Error() string {
	return err.Message
}

type Client struct {
	c       *http.Client
	baseURL *url.URL
}

func (c *Client) SetBaseURL(base string) error {
	var err error
	c.baseURL, err = url.Parse(base)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetHTTPClient(hc *http.Client) error {
	c.c = hc
	return nil
}

func (c *Client) Login(ctx context.Context, user, pass string) error {
	if c.c == nil {
		c.c = &http.Client{}

		jar, _ := cookiejar.New(nil)
		c.c.Jar = jar
	}

	err := c.doLogin(ctx, "POST", "api/auth/login", &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: user,
		Password: pass,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doLogin(ctx context.Context, method, relativeURL string, reqBody interface{}, respBody interface{}) error {
	var (
		reqReader io.Reader
		err       error
		reqBytes  []byte
	)
	if reqBody != nil {

		reqBytes, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("unable to marshal JSON: %s %s %w", method, relativeURL, err)
		}
		reqReader = bytes.NewReader(reqBytes)
	}

	reqURL, err := url.Parse(relativeURL)
	if err != nil {
		return fmt.Errorf("unable to parse URL: %s %s %w", method, relativeURL, err)
	}

	url := c.baseURL.ResolveReference(reqURL)

	req, err := http.NewRequestWithContext(ctx, method, url.String(), reqReader)
	if err != nil {
		return fmt.Errorf("unable to create request: %s %s %w", method, relativeURL, err)
	}

	req.Header.Set("User-Agent", "terraform-provider-unifi/0.1")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s %s %w", method, relativeURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{}
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Request Body:\n%s\n", string(reqBytes))
		errBody := struct {
			Meta meta `json:"meta"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("%s %s (%s) for %s %s", errBody.Meta.RC, errBody.Meta.Message, resp.Status, method, url.String())
	}

	if respBody == nil || resp.ContentLength == 0 {
		return nil
	}

	// TODO: check rc in addition to status code?

	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		return fmt.Errorf("unable to decode body: %s %s %w", method, relativeURL, err)
	}

	return nil
}

func (c *Client) do(ctx context.Context, method, relativeURL string, reqBody interface{}, respBody interface{}) error {
	var (
		reqReader io.Reader
		err       error
		reqBytes  []byte
	)
	if reqBody != nil {

		reqBytes, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("unable to marshal JSON: %s %s %w", method, relativeURL, err)
		}
		reqReader = bytes.NewReader(reqBytes)
	}

	reqURL, err := url.Parse(relativeURL)
	if err != nil {
		return fmt.Errorf("unable to parse URL: %s %s %w", method, relativeURL, err)
	}

	url := c.baseURL.ResolveReference(reqURL)

	req, err := http.NewRequestWithContext(ctx, method, url.String(), reqReader)
	if err != nil {
		return fmt.Errorf("unable to create request: %s %s %w", method, relativeURL, err)
	}

	req.Header.Set("User-Agent", "terraform-provider-unifi/0.1")

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s %s %w", method, relativeURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{}
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Request Body:\n%s\n", string(reqBytes))
		errBody := struct {
			Meta meta `json:"meta"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("%s %s (%s) for %s %s", errBody.Meta.RC, errBody.Meta.Message, resp.Status, method, url.String())
	}

	if respBody == nil || resp.ContentLength == 0 {
		return nil
	}

	// TODO: check rc in addition to status code?

	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		return fmt.Errorf("unable to decode body: %s %s %w", method, relativeURL, err)
	}

	return nil
}

type meta struct {
	RC      string `json:"rc"`
	Message string `json:"msg"`
}

func (m *meta) error() error {
	if m.RC != "ok" {
		return &APIError{
			RC:      m.RC,
			Message: m.Message,
		}
	}

	return nil
}

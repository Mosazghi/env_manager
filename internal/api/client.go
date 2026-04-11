package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

func NewClient(token string, baseURL string) *Client {
	return &Client{token: token, httpClient: &http.Client{}, baseURL: baseURL}
}

func (c *Client) do(method, path string, body any) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	}
	apiPath := fmt.Sprintf("%v/api%v", c.baseURL, path)

	if !isValidPath(apiPath) {
		return nil, fmt.Errorf("invalid path: %v", apiPath)
	}

	req, err := http.NewRequest(method, apiPath, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	parsedBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldnt parse response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error: %v", string(parsedBody))
	}

	return parsedBody, nil
}

func (c *Client) Get(path string) ([]byte, error)            { return c.do("GET", path, nil) }
func (c *Client) Post(path string, body any) ([]byte, error) { return c.do("POST", path, body) }
func (c *Client) Delete(path string) ([]byte, error)         { return c.do("DELETE", path, nil) }
func (c *Client) Put(path string, body any) ([]byte, error)  { return c.do("PUT", path, body) }

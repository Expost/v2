package mercury

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	URL string
}

type MercuryResp struct {
	Title         string
	Author        string
	DatePublished string
	URL           string
	Domain        string
	WordCount     int
	Direction     string
	TotalPages    int
	RenderedPages int
	Content       string
}

func NewClient(url string) *Client {
	return &Client{URL: url}
}

func (c *Client) GetFullText(entryURL string) (string, error) {
	if c.URL == "" {
		return "", fmt.Errorf("mercury: missing mercury api url")
	}

	reqURL, err := url.Parse(c.URL)
	if err != nil {
		return "", fmt.Errorf("mercury: failed to parser url, %v", err)
	}

	queryParams := url.Values{}
	queryParams.Set("url", entryURL)

	reqURL.Path = "/parser"
	reqURL.RawQuery = queryParams.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		fmt.Println("Failed to send GET request:", err)
		return "", fmt.Errorf("mercury: failed to send GET request, %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("mercury: unable to get fulltext: url=%s status=%d", reqURL.String(), resp.StatusCode)
	}

	var mercuryRsp MercuryResp
	if err := json.NewDecoder(resp.Body).Decode(&mercuryRsp); err != nil {
		return "", fmt.Errorf("mercury: uanble to decode mercury response: %v", err)
	}

	return mercuryRsp.Content, nil
}

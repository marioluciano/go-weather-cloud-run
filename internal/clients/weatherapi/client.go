package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client
	apiKey string
}

func NewClient(client *http.Client, apiKey string) *Client {
	return &Client{client, apiKey}
}

func (c *Client) FetchByCity(ctx context.Context, city string) (*WeatherAPIResponse, error) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", c.apiKey, url.QueryEscape(city))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response WeatherAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

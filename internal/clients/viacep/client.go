package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{client}
}

func (c *Client) FetchByCep(ctx context.Context, cep string) (*ViaCEPResponse, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

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

	var response ViaCEPResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

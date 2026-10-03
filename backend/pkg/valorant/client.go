package valorant

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiBase = "https://valorant-api.com/v1"

type Client struct {
	HTTP *http.Client
}

func NewClient() *Client {
	return &Client{HTTP: &http.Client{Timeout: 10 * time.Second}}
}

type Map struct {
	UUID                string `json:"uuid"`
	DisplayName         string `json:"displayName"`
	TacticalDescription string `json:"tacticalDescription"`
}

func (c *Client) Maps(ctx context.Context) ([]Map, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBase+"/maps", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("maps: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("maps: status %d", res.StatusCode)
	}

	var body struct {
		Data []Map `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("maps: %w", err)
	}
	return body.Data, nil
}

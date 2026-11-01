package valorant

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const apiBase = "https://valorant-api.com/v1"

// Os splashes têm uns 2,3 MB; o limite só evita guardar algo absurdo se a API mudar.
const maxImageBytes = 8 << 20

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
	ListViewIcon        string `json:"listViewIcon"`
	Splash              string `json:"splash"`
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

// Image baixa uma arte de media.valorant-api.com, que o Discord não deixa a atividade abrir direto.
func (c *Client) Image(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("image: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("image: status %d", res.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(res.Body, maxImageBytes))
	if err != nil {
		return nil, "", fmt.Errorf("image: %w", err)
	}
	return data, res.Header.Get("Content-Type"), nil
}

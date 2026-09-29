// Package discord faz as chamadas REST ao Discord usadas pela Activity.
package discord

import (
	"context"
	"errors"
)

type Client struct {
	ClientID     string
	ClientSecret string
}

type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Avatar     string `json:"avatar"`
}

// ExchangeCode troca o code do SDK por um access token.
func (c *Client) ExchangeCode(ctx context.Context, code string) (string, error) {
	return "", errors.New("não implementado")
}

// CurrentUser devolve o usuário dono do access token.
func (c *Client) CurrentUser(ctx context.Context, accessToken string) (*User, error) {
	return nil, errors.New("não implementado")
}

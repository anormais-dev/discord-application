package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const apiBase = "https://discord.com/api/v10"

type Client struct {
	ClientID     string
	ClientSecret string
	HTTP         *http.Client
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTP:         &http.Client{Timeout: 10 * time.Second},
	}
}

type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Avatar     string `json:"avatar"`
	Nick       string `json:"-"`
}

type Member struct {
	Nick string `json:"nick"`
}

func (u *User) DisplayName() string {
	if u.GlobalName != "" {
		return u.GlobalName
	}
	return u.Username
}

func (u *User) AvatarURL() string {
	if u.Avatar != "" {
		return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=128", u.ID, u.Avatar)
	}
	id, _ := strconv.ParseUint(u.ID, 10, 64)
	return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", (id>>22)%6)
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (string, error) {
	form := url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBase+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var body struct {
		AccessToken string `json:"access_token"`
	}
	if err := c.do(req, &body); err != nil {
		return "", fmt.Errorf("oauth2 token: %w", err)
	}
	return body.AccessToken, nil
}

func (c *Client) CurrentUser(ctx context.Context, accessToken string) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBase+"/users/@me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var u User
	if err := c.do(req, &u); err != nil {
		return nil, fmt.Errorf("users/@me: %w", err)
	}
	return &u, nil
}

// GuildMember precisa do escopo guilds.members.read.
func (c *Client) GuildMember(ctx context.Context, accessToken, guildID string) (*Member, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBase+"/users/@me/guilds/"+url.PathEscape(guildID)+"/member", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var m Member
	if err := c.do(req, &m); err != nil {
		return nil, fmt.Errorf("guild member: %w", err)
	}
	return &m, nil
}

func (c *Client) do(req *http.Request, out any) error {
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

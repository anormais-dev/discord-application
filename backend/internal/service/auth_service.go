package service

import (
	"context"
	"errors"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

// DevTokenPrefix marca um token de desenvolvimento: "dev:<nome>".
const DevTokenPrefix = "dev:"

const maxDevNameLen = 32

var (
	ErrNoDiscordClient = errors.New("credenciais do Discord não configuradas")
	ErrDevNameRequired = errors.New("nome obrigatório")
)

type DiscordClient interface {
	ExchangeCode(ctx context.Context, code string) (string, error)
	CurrentUser(ctx context.Context, accessToken string) (*discord.User, error)
}

type AuthService struct {
	discord DiscordClient
	devAuth bool
}

// NewAuthService aceita discord nil quando não há credenciais (só faz sentido com devAuth).
func NewAuthService(discord DiscordClient, devAuth bool) *AuthService {
	return &AuthService{discord: discord, devAuth: devAuth}
}

func (s *AuthService) ExchangeCode(ctx context.Context, code string) (string, error) {
	if s.discord == nil {
		return "", ErrNoDiscordClient
	}
	return s.discord.ExchangeCode(ctx, code)
}

func (s *AuthService) CurrentUser(ctx context.Context, accessToken string) (*discord.User, error) {
	if s.devAuth {
		if name, ok := strings.CutPrefix(accessToken, DevTokenPrefix); ok {
			return devUser(name)
		}
	}
	if s.discord == nil {
		return nil, ErrNoDiscordClient
	}
	return s.discord.CurrentUser(ctx, accessToken)
}

func devUser(name string) (*discord.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrDevNameRequired
	}
	if r := []rune(name); len(r) > maxDevNameLen {
		name = string(r[:maxDevNameLen])
	}
	// ID numérico estável por nome, para o avatar padrão variar entre usuários.
	h := fnv.New64a()
	h.Write([]byte(strings.ToLower(name)))
	id := strconv.FormatUint(h.Sum64()>>1, 10)
	return &discord.User{ID: id, Username: name, GlobalName: name}, nil
}

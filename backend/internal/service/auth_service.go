package service

import (
	"context"
	"errors"
	"hash/fnv"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

// DevTokenPrefix marca um token de desenvolvimento: "dev:<nome>".
const DevTokenPrefix = "dev:"

const maxDevNameLen = 32

// Quanto tempo o usuário de um token fica guardado. Cada reconexão do WebSocket
// autentica de novo, e sem cache isso bate no Discord toda vez (e esbarra no rate limit).
const userCacheTTL = 10 * time.Minute

var (
	ErrNoDiscordClient = errors.New("credenciais do Discord não configuradas")
	ErrDevNameRequired = errors.New("nome obrigatório")
)

type DiscordClient interface {
	ExchangeCode(ctx context.Context, code string) (string, error)
	CurrentUser(ctx context.Context, accessToken string) (*discord.User, error)
	GuildMember(ctx context.Context, accessToken, guildID string) (*discord.Member, error)
}

type AuthService struct {
	discord DiscordClient
	devAuth bool

	mu    sync.Mutex
	users map[userCacheKey]cachedUser
}

type userCacheKey struct {
	accessToken string
	guildID     string
}

type cachedUser struct {
	user    discord.User
	expires time.Time
}

// NewAuthService aceita discord nil quando não há credenciais (só faz sentido com devAuth).
func NewAuthService(discord DiscordClient, devAuth bool) *AuthService {
	return &AuthService{discord: discord, devAuth: devAuth, users: map[userCacheKey]cachedUser{}}
}

func (s *AuthService) ExchangeCode(ctx context.Context, code string) (string, error) {
	if s.discord == nil {
		return "", ErrNoDiscordClient
	}
	return s.discord.ExchangeCode(ctx, code)
}

func (s *AuthService) CurrentUser(ctx context.Context, accessToken, guildID string) (*discord.User, error) {
	if s.devAuth {
		if name, ok := strings.CutPrefix(accessToken, DevTokenPrefix); ok {
			return devUser(name)
		}
	}
	if s.discord == nil {
		return nil, ErrNoDiscordClient
	}
	key := userCacheKey{accessToken, guildID}
	if u, ok := s.cached(key); ok {
		return u, nil
	}

	u, err := s.discord.CurrentUser(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if guildID != "" {
		if m, err := s.discord.GuildMember(ctx, accessToken, guildID); err != nil {
			log.Printf("apelido no servidor: %v", err)
		} else {
			u.Nick = m.Nick
		}
	}
	s.store(key, u)
	return u, nil
}

func (s *AuthService) cached(key userCacheKey) (*discord.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.users[key]
	if !ok || time.Now().After(e.expires) {
		return nil, false
	}
	u := e.user
	return &u, true
}

// store aproveita para jogar fora o que já venceu; o mapa nunca passa de quem usou o app nos últimos minutos.
func (s *AuthService) store(key userCacheKey, u *discord.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, e := range s.users {
		if now.After(e.expires) {
			delete(s.users, k)
		}
	}
	s.users[key] = cachedUser{user: *u, expires: now.Add(userCacheTTL)}
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

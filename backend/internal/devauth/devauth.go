// Package devauth permite testar a roleta fora do Discord, com usuários falsos.
// Nunca deve ser ativado em produção: qualquer um consegue se passar por qualquer nome.
package devauth

import (
	"context"
	"errors"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

// TokenPrefix marca um token de desenvolvimento: "dev:<nome>".
const TokenPrefix = "dev:"

const maxNameLen = 32

// Users aceita tokens "dev:<nome>" e repassa os demais para o lookup real, se houver.
type Users struct {
	Next interface {
		CurrentUser(ctx context.Context, accessToken string) (*discord.User, error)
	}
}

func (u Users) CurrentUser(ctx context.Context, token string) (*discord.User, error) {
	name, ok := strings.CutPrefix(token, TokenPrefix)
	if !ok {
		if u.Next == nil {
			return nil, errors.New("token de desenvolvimento inválido")
		}
		return u.Next.CurrentUser(ctx, token)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("nome obrigatório")
	}
	if r := []rune(name); len(r) > maxNameLen {
		name = string(r[:maxNameLen])
	}
	// ID numérico estável por nome, para o avatar padrão variar entre usuários.
	h := fnv.New64a()
	h.Write([]byte(strings.ToLower(name)))
	id := strconv.FormatUint(h.Sum64()>>1, 10)
	return &discord.User{ID: id, Username: name, GlobalName: name}, nil
}

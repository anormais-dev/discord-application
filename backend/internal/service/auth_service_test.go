package service

import (
	"context"
	"errors"
	"testing"

	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

func TestDevToken(t *testing.T) {
	s := NewAuthService(nil, true)
	a, err := s.CurrentUser(context.Background(), "dev:Ana", "")
	if err != nil || a.DisplayName() != "Ana" {
		t.Fatalf("esperado usuário Ana, veio %+v %v", a, err)
	}
	again, _ := s.CurrentUser(context.Background(), "dev:ana", "")
	if again.ID != a.ID {
		t.Fatalf("o mesmo nome deveria gerar o mesmo ID")
	}
	b, _ := s.CurrentUser(context.Background(), "dev:Bruno", "")
	if b.ID == a.ID {
		t.Fatalf("nomes diferentes deveriam gerar IDs diferentes")
	}
}

func TestRejectsWithoutDiscordClient(t *testing.T) {
	s := NewAuthService(nil, true)
	for _, token := range []string{"dev:  ", "token-real"} {
		if _, err := s.CurrentUser(context.Background(), token, ""); err == nil {
			t.Fatalf("token %q deveria ser recusado", token)
		}
	}
}

func TestDevTokenIgnoredWithoutDevAuth(t *testing.T) {
	s := NewAuthService(nil, false)
	if _, err := s.CurrentUser(context.Background(), "dev:Ana", ""); err == nil {
		t.Fatalf("token de dev não deveria valer sem DEV_AUTH")
	}
}

type nickDiscord struct {
	memberErr error
}

func (nickDiscord) ExchangeCode(context.Context, string) (string, error) { return "", nil }

func (nickDiscord) CurrentUser(context.Context, string) (*discord.User, error) {
	return &discord.User{ID: "1", Username: "ana", GlobalName: "Ana Global"}, nil
}

func (d nickDiscord) GuildMember(context.Context, string, string) (*discord.Member, error) {
	if d.memberErr != nil {
		return nil, d.memberErr
	}
	return &discord.Member{Nick: "Aninha"}, nil
}

func TestCurrentUserUsesGuildNick(t *testing.T) {
	u, err := NewAuthService(nickDiscord{}, false).CurrentUser(context.Background(), "t", "g1")
	if err != nil || u.Nick != "Aninha" {
		t.Fatalf("esperado apelido Aninha, veio %+v %v", u, err)
	}
	u, _ = NewAuthService(nickDiscord{}, false).CurrentUser(context.Background(), "t", "")
	if u.Nick != "" {
		t.Fatalf("sem guild não deveria buscar apelido, veio %q", u.Nick)
	}
}

func TestCurrentUserIgnoresGuildMemberError(t *testing.T) {
	u, err := NewAuthService(nickDiscord{memberErr: errors.New("403")}, false).CurrentUser(context.Background(), "t", "g1")
	if err != nil || u.Nick != "" || u.GlobalName != "Ana Global" {
		t.Fatalf("falha no apelido não deveria impedir o login: %+v %v", u, err)
	}
}

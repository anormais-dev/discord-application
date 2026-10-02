package service

import (
	"context"
	"testing"
)

func TestDevToken(t *testing.T) {
	s := NewAuthService(nil, true)
	a, err := s.CurrentUser(context.Background(), "dev:Ana")
	if err != nil || a.DisplayName() != "Ana" {
		t.Fatalf("esperado usuário Ana, veio %+v %v", a, err)
	}
	again, _ := s.CurrentUser(context.Background(), "dev:ana")
	if again.ID != a.ID {
		t.Fatalf("o mesmo nome deveria gerar o mesmo ID")
	}
	b, _ := s.CurrentUser(context.Background(), "dev:Bruno")
	if b.ID == a.ID {
		t.Fatalf("nomes diferentes deveriam gerar IDs diferentes")
	}
}

func TestRejectsWithoutDiscordClient(t *testing.T) {
	s := NewAuthService(nil, true)
	for _, token := range []string{"dev:  ", "token-real"} {
		if _, err := s.CurrentUser(context.Background(), token); err == nil {
			t.Fatalf("token %q deveria ser recusado", token)
		}
	}
}

func TestDevTokenIgnoredWithoutDevAuth(t *testing.T) {
	s := NewAuthService(nil, false)
	if _, err := s.CurrentUser(context.Background(), "dev:Ana"); err == nil {
		t.Fatalf("token de dev não deveria valer sem DEV_AUTH")
	}
}

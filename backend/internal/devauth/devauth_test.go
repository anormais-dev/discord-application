package devauth

import (
	"context"
	"testing"
)

func TestDevToken(t *testing.T) {
	u := Users{}
	a, err := u.CurrentUser(context.Background(), "dev:Ana")
	if err != nil || a.DisplayName() != "Ana" {
		t.Fatalf("esperado usuário Ana, veio %+v %v", a, err)
	}
	again, _ := u.CurrentUser(context.Background(), "dev:ana")
	if again.ID != a.ID {
		t.Fatalf("o mesmo nome deveria gerar o mesmo ID")
	}
	b, _ := u.CurrentUser(context.Background(), "dev:Bruno")
	if b.ID == a.ID {
		t.Fatalf("nomes diferentes deveriam gerar IDs diferentes")
	}
}

func TestRejectsWithoutNext(t *testing.T) {
	u := Users{}
	for _, token := range []string{"dev:  ", "token-real"} {
		if _, err := u.CurrentUser(context.Background(), token); err == nil {
			t.Fatalf("token %q deveria ser recusado", token)
		}
	}
}

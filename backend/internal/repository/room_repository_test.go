package repository

import (
	"errors"
	"testing"
)

func TestRoomRepository(t *testing.T) {
	repo := NewRoomRepository()

	if _, err := repo.Get("abc"); !errors.Is(err, ErrRoomNotFound) {
		t.Fatalf("erro = %v, esperado ErrRoomNotFound", err)
	}

	created := repo.GetOrCreate("abc")
	if again := repo.GetOrCreate("abc"); again != created {
		t.Fatalf("GetOrCreate deveria devolver a mesma room para a mesma instância")
	}
	if got, err := repo.Get("abc"); err != nil || got != created {
		t.Fatalf("Get = %v %v, esperado a room criada", got, err)
	}

	repo.Delete("abc")
	if _, err := repo.Get("abc"); !errors.Is(err, ErrRoomNotFound) {
		t.Fatalf("room deveria ter sido removida")
	}
}

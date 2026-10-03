package service

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/anormais-dev/discord-application/backend/pkg/valorant"
)

type fakeValorant struct {
	maps  []valorant.Map
	err   error
	calls int
}

func (f *fakeValorant) Maps(context.Context) ([]valorant.Map, error) {
	f.calls++
	return f.maps, f.err
}

var apiMaps = []valorant.Map{
	{DisplayName: "Haven", TacticalDescription: "A/B/C Sites"},
	{DisplayName: "The Range"},
	{DisplayName: "Ascent", TacticalDescription: "A/B Sites"},
	{DisplayName: "Skirmish A"},
	{DisplayName: "Ascent", TacticalDescription: "A/B Sites"},
}

func TestCompetitiveMapsFiltersAndSorts(t *testing.T) {
	s := NewMapService(&fakeValorant{maps: apiMaps})
	got, err := s.CompetitiveMaps(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"Ascent", "Haven"}; !slices.Equal(got, want) {
		t.Fatalf("mapas = %v, esperado %v", got, want)
	}
}

func TestCompetitiveMapsUsesCacheAndFallsBackWhenAPIFails(t *testing.T) {
	api := &fakeValorant{maps: apiMaps}
	now := time.Now()
	s := NewMapService(api)
	s.now = func() time.Time { return now }

	s.CompetitiveMaps(context.Background())
	s.CompetitiveMaps(context.Background())
	if api.calls != 1 {
		t.Fatalf("dentro do TTL deveria usar o cache, chamou a API %d vezes", api.calls)
	}

	now = now.Add(mapsTTL + time.Minute)
	api.err = errors.New("fora do ar")
	got, err := s.CompetitiveMaps(context.Background())
	if err != nil || len(got) != 2 || api.calls != 2 {
		t.Fatalf("com a API fora deveria devolver o cache antigo, veio %v, %v (chamadas: %d)", got, err, api.calls)
	}
}

func TestCompetitiveMapsWithoutCacheFails(t *testing.T) {
	s := NewMapService(&fakeValorant{err: errors.New("fora do ar")})
	if _, err := s.CompetitiveMaps(context.Background()); !errors.Is(err, ErrMapsUnavailable) {
		t.Fatalf("erro = %v, esperado ErrMapsUnavailable", err)
	}
}

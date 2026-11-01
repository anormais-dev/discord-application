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
	maps      []valorant.Map
	err       error
	calls     int
	imageURLs []string
}

func (f *fakeValorant) Maps(context.Context) ([]valorant.Map, error) {
	f.calls++
	return f.maps, f.err
}

func (f *fakeValorant) Image(_ context.Context, url string) ([]byte, string, error) {
	f.imageURLs = append(f.imageURLs, url)
	return []byte("png de " + url), "image/png", nil
}

var apiMaps = []valorant.Map{
	{UUID: "haven-id", DisplayName: "Haven", TacticalDescription: "A/B/C Sites", ListViewIcon: "haven-card.png", Splash: "haven-splash.png"},
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

func TestImagesComeWithTheCompetitiveMaps(t *testing.T) {
	s := NewMapService(&fakeValorant{maps: apiMaps})
	if _, err := s.CompetitiveMaps(context.Background()); err != nil {
		t.Fatal(err)
	}
	images := s.Images()
	if got := images["Haven"]; got != (MapImages{UUID: "haven-id", Card: "haven-card.png", Splash: "haven-splash.png"}) {
		t.Fatalf("imagens do Haven = %+v", got)
	}
	if _, ok := images["The Range"]; ok {
		t.Fatalf("mapa fora do competitivo não deveria ter imagem")
	}
}

func TestImageDownloadsOnceAndOnlyKnownMaps(t *testing.T) {
	api := &fakeValorant{maps: apiMaps}
	s := NewMapService(api)

	for range 2 {
		img, err := s.Image(context.Background(), "haven-id", "splash")
		if err != nil {
			t.Fatal(err)
		}
		if string(img.Data) != "png de haven-splash.png" || img.ContentType != "image/png" {
			t.Fatalf("imagem = %+v", img)
		}
	}
	if !slices.Equal(api.imageURLs, []string{"haven-splash.png"}) {
		t.Fatalf("deveria baixar uma vez só, baixou %v", api.imageURLs)
	}

	for _, c := range [][2]string{{"outro-id", "card"}, {"haven-id", "minimap"}} {
		if _, err := s.Image(context.Background(), c[0], c[1]); !errors.Is(err, ErrMapImageNotFound) {
			t.Fatalf("%v: erro = %v, esperado ErrMapImageNotFound", c, err)
		}
	}
}

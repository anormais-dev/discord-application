package service

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/anormais-dev/discord-application/backend/pkg/valorant"
)

// Os mapas mudam só com patch do jogo; meio dia de cache basta.
const mapsTTL = 12 * time.Hour

type ValorantClient interface {
	Maps(ctx context.Context) ([]valorant.Map, error)
	Image(ctx context.Context, url string) ([]byte, string, error)
}

type MapService struct {
	client ValorantClient
	now    func() time.Time

	mu        sync.Mutex
	maps      []string
	images    map[string]MapImages
	fetchedAt time.Time

	// Artes já baixadas, pela URL de origem. Mudam só com patch do jogo, então não expiram.
	imgMu  sync.Mutex
	imgSet map[string]MapImage
}

// MapImages são as artes de um mapa na Valorant API: a faixa leve para o carrossel e o
// splash grande para o resultado. O UUID identifica o mapa na rota que repassa as imagens.
type MapImages struct {
	UUID   string
	Card   string
	Splash string
}

// MapImage é uma arte baixada, pronta para devolver ao navegador.
type MapImage struct {
	Data        []byte
	ContentType string
}

func NewMapService(client ValorantClient) *MapService {
	return &MapService{client: client, now: time.Now}
}

// CompetitiveMaps devolve os nomes dos mapas com bomb, que são os únicos com
// descrição tática na API. Se a API cair, segue com a última lista que deu certo.
func (s *MapService) CompetitiveMaps(ctx context.Context) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.maps != nil && s.now().Sub(s.fetchedAt) < mapsTTL {
		return s.maps, nil
	}
	all, err := s.client.Maps(ctx)
	names := competitive(all)
	if err != nil || len(names) == 0 {
		if s.maps != nil {
			return s.maps, nil
		}
		return nil, ErrMapsUnavailable
	}
	s.maps = names
	s.images = imagesByName(all)
	s.fetchedAt = s.now()
	return s.maps, nil
}

func competitive(all []valorant.Map) []string {
	names := []string{}
	for _, m := range all {
		if m.TacticalDescription != "" && !slices.Contains(names, m.DisplayName) {
			names = append(names, m.DisplayName)
		}
	}
	slices.Sort(names)
	return names
}

// Images devolve as artes da última lista buscada em CompetitiveMaps.
func (s *MapService) Images() map[string]MapImages {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.images
}

func imagesByName(all []valorant.Map) map[string]MapImages {
	images := map[string]MapImages{}
	for _, m := range all {
		if m.TacticalDescription != "" {
			images[m.DisplayName] = MapImages{UUID: m.UUID, Card: m.ListViewIcon, Splash: m.Splash}
		}
	}
	return images
}

// Image baixa (uma vez) a arte card ou splash de um mapa competitivo. Só serve URLs que
// vieram da própria API, para a rota não virar um proxy aberto.
func (s *MapService) Image(ctx context.Context, uuid, kind string) (MapImage, error) {
	// Garante a lista carregada: depois de reiniciar, o navegador pode pedir a imagem antes da lista.
	if _, err := s.CompetitiveMaps(ctx); err != nil {
		return MapImage{}, err
	}
	url := s.imageURL(uuid, kind)
	if url == "" {
		return MapImage{}, ErrMapImageNotFound
	}

	s.imgMu.Lock()
	img, ok := s.imgSet[url]
	s.imgMu.Unlock()
	if ok {
		return img, nil
	}

	data, contentType, err := s.client.Image(ctx, url)
	if err != nil {
		return MapImage{}, ErrMapsUnavailable
	}
	img = shrink(data, contentType)
	s.imgMu.Lock()
	if s.imgSet == nil {
		s.imgSet = map[string]MapImage{}
	}
	s.imgSet[url] = img
	s.imgMu.Unlock()
	return img, nil
}

func (s *MapService) imageURL(uuid, kind string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, img := range s.images {
		if img.UUID != uuid {
			continue
		}
		switch kind {
		case "card":
			return img.Card
		case "splash":
			return img.Splash
		}
	}
	return ""
}

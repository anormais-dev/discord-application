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
}

type MapService struct {
	client ValorantClient
	now    func() time.Time

	mu        sync.Mutex
	maps      []string
	fetchedAt time.Time
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

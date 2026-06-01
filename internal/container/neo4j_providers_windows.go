//go:build windows

package container

import (
	"context"

	"go.uber.org/dig"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// registerNeo4jProviders registers no-op stubs for neo4j-related providers.
// Windows builds exclude the neo4j-go-driver to work around a Go 1.26.3
// compiler bug ("internal error: package without types was imported").
func registerNeo4jProviders(c *dig.Container) {
	must(c.Provide(func() interfaces.RetrieveGraphRepository { return &stubGraphRepo{} }))
	must(c.Provide(func() interfaces.MemoryRepository { return &stubMemoryRepo{} }))
	must(c.Provide(func() string { return "" })) // graphEngine: not available
}

// ── stub implementations ──

type stubGraphRepo struct{}

func (s *stubGraphRepo) AddGraph(_ context.Context, _ types.NameSpace, _ []*types.GraphData) error {
	return nil
}

func (s *stubGraphRepo) DelGraph(_ context.Context, _ []types.NameSpace) error {
	return nil
}

func (s *stubGraphRepo) SearchNode(_ context.Context, _ types.NameSpace, _ []string) (*types.GraphData, error) {
	return nil, nil
}

type stubMemoryRepo struct{}

func (s *stubMemoryRepo) SaveEpisode(_ context.Context, _ *types.Episode, _ []*types.Entity, _ []*types.Relationship) error {
	return nil
}

func (s *stubMemoryRepo) FindRelatedEpisodes(_ context.Context, _ string, _ []string, _ int) ([]*types.Episode, error) {
	return nil, nil
}

func (s *stubMemoryRepo) IsAvailable(_ context.Context) bool {
	return false
}

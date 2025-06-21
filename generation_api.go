package pokesdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jameshalsall/pokesdk/internal/backend"
	"github.com/jameshalsall/pokesdk/internal/urlutil"
)

const (
	apiGenerationPath = "/generation"
)

type GenerationAPI struct {
	cfg Config
}

// List returns a Paginator for listing all Generations.
// It accepts no context argument because it should be provided to the paginator's functions instead.
func (g GenerationAPI) List() *Paginator[*GenerationList] {
	response := &GenerationList{}

	return NewPaginator[*GenerationList](g.url(apiGenerationPath), func(ctx context.Context, nextUrl string) (*GenerationList, error) {
		err := g.cfg.backend.Process(ctx, nextUrl, nil, response)
		if err != nil {
			return nil, fmt.Errorf("pokesdk: error listing generations: %w", err)
		}

		return response, nil
	})
}

// GetByName retrieves a specific Generation by its name.
func (g GenerationAPI) GetByName(ctx context.Context, name string) (*Generation, error) {
	return g.getGeneration(ctx, g.url(apiGenerationPath+"/"+name))
}

// GetByID retrieves a specific Generation by its ID.
func (g GenerationAPI) GetByID(ctx context.Context, ID int) (*Generation, error) {
	return g.getGeneration(ctx, g.url(apiGenerationPath+"/"+strconv.Itoa(ID)))
}

// GetByRef retrieves a specific Generation by its reference.
// The reference is returned in the response from List()
func (g GenerationAPI) GetByRef(ctx context.Context, ref GenerationRef) (*Generation, error) {
	return g.getGeneration(ctx, ref.URL)
}

func (g GenerationAPI) getGeneration(ctx context.Context, url string) (*Generation, error) {
	response := &Generation{}
	err := g.cfg.backend.Process(ctx, url, nil, response)
	if err != nil {
		if errors.Is(err, backend.ErrResourceNotFound) {
			return nil, ErrGenerationNotFound
		}
		return nil, fmt.Errorf("pokesdk: error getting generation: %w", err)
	}

	return response, nil
}

func (g GenerationAPI) url(path string) string {
	return urlutil.BuildURL(g.cfg.baseURL, path)
}

package pokesdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jameshalsall/pokesdk/internal/backend"
	oururl "github.com/jameshalsall/pokesdk/internal/url"
)

const (
	apiPokemonPath = "/pokemon"
)

type PokemonAPI struct {
	cfg Config
}

// List returns a Paginator for listing all Pokemon.
// It accepts no context argument because it should be provided to the paginator's functions instead.
func (g PokemonAPI) List() *Paginator[*PokemonList] {
	response := &PokemonList{}

	return NewPaginator[*PokemonList](g.url(apiPokemonPath), func(ctx context.Context, nextUrl string) (*PokemonList, error) {
		err := g.cfg.backend.Process(ctx, nextUrl, nil, response)
		if err != nil {
			return nil, fmt.Errorf("pokesdk: error listing pokemon: %w", err)
		}

		return response, nil
	})
}

// GetByName retrieves a specific Pokemon by its name.
func (g PokemonAPI) GetByName(ctx context.Context, name string) (*Pokemon, error) {
	return g.getPokemon(ctx, g.url(apiPokemonPath+"/"+name))
}

// GetByID retrieves a specific Pokemon by its ID.
func (g PokemonAPI) GetByID(ctx context.Context, ID int) (*Pokemon, error) {
	return g.getPokemon(ctx, g.url(apiPokemonPath+"/"+strconv.Itoa(ID)))
}

// GetByRef retrieves a specific Pokemon by its reference.
// The reference is returned in the response from List()
func (g PokemonAPI) GetByRef(ctx context.Context, ref PokemonRef) (*Pokemon, error) {
	return g.getPokemon(ctx, ref.URL)
}

func (g PokemonAPI) getPokemon(ctx context.Context, url string) (*Pokemon, error) {
	response := &Pokemon{}
	err := g.cfg.backend.Process(ctx, url, nil, response)
	if err != nil {
		if errors.Is(err, backend.ErrResourceNotFound) {
			return nil, ErrPokemonNotFound
		}
		return nil, fmt.Errorf("pokesdk: error getting pokemon: %w", err)
	}

	return response, nil
}

func (g PokemonAPI) url(path string) string {
	return oururl.BuildURL(g.cfg.baseURL, path)
}

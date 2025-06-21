package pokesdk

import (
	"errors"
	"fmt"
)

var (
	errNotFound           = errors.New("not found")
	ErrPokemonNotFound    = fmt.Errorf("pokemon: %w", errNotFound)
	ErrGenerationNotFound = fmt.Errorf("generation: %w", errNotFound)
)

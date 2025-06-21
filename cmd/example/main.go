package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jameshalsall/pokesdk"
)

// Example of using the pokesdk to fetch a list of Pokémon and their details.

func main() {
	client := pokesdk.NewClient()

	pokemonList := client.Pokemon.List()
	firstPage := pokemonList.Next(context.Background())
	if firstPage.Error != nil {
		slog.Error("Error fetching Pokémon list", "error", firstPage.Error)
		os.Exit(1)
	}

	var lastRef pokesdk.PokemonRef
	for _, pokemon := range firstPage.Result.Results {
		slog.Info("Pokémon found", "name", pokemon.Name)
		lastRef = pokemon
	}

	detailedPokemon, err := client.Pokemon.GetByRef(context.Background(), lastRef)
	if err != nil {
		slog.Error("Error fetching Pokémon details", "error", err)
		os.Exit(1)
	}
	slog.Info("Detailed Pokémon information", "name", detailedPokemon.Name, "id", detailedPokemon.ID, "height", detailedPokemon.Height, "weight", detailedPokemon.Weight)

	// fetch the next page of Pokémon
	secondPage := pokemonList.Next(context.Background())
	if pokemonList == nil {
		slog.Info("no more Pokémon")
		return
	}

	for _, pokemon := range secondPage.Result.Results {
		slog.Info("Pokémon found", "name", pokemon.Name)
	}
}

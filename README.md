# Poke API SDK

SDK for V2 of the [PokéAPI](https://pokeapi.co), a RESTful API for accessing data about Pokémon.

## Installation

```bash
go get github.com/jameshalsall/poke-sdk
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jameshalsall/pokesdk"
)

func main() {
	// Create a new client
	client := pokesdk.NewClient()

	pikachu, err := client.Pokemon.GetByName(context.Background(), "pikachu")
	if err != nil {
		fmt.Printf("Error fetching Pokémon: %v", err)
		os.Exit(1)
	}

	slog.Info("Fetched Pokémon: %s", pikachu.Name)
}
```

## How to use

The API requires no authentication, so you can start using it immediately.

### Configuring the client
The client can be configured with various functional options.

```go
myClient := &http.Client{}

client := pokesdk.NewClient(
	pokesdk.WithCustomHttpClient(myClient),
	pokesdk.WithCustomBaseURL("https://custom.pokeapi.co/api/v2/"),
)
```

## Getting Pokemon data
### Listing all Pokémon

```go
client := pokesdk.NewClient()

// Fetch first page of Pokémon (use a better context in production code)
pokemonList := client.Pokemon.List()
firstPage := pokemonList.Next(context.Background())
if firstPage.Error != nil {
    slog.Error("Error fetching Pokémon list", "error", firstPage.Error)
	os.Exit(1)
}

for _, pokemon := range firstPage.Result.Results {
    slog.Info("Pokémon found", "name", pokemon.Name)
	// you can get more information on the pokemon by fetching it by ref
	detailedPokemon, err := client.Pokemon.GetByRef(context.Background(), pokemon)
	if err != nil {
        slog.Error("Error fetching Pokémon details", "error", err)
		continue
	}
}

// fetch the next page of Pokémon
secondPage := pokemonList.Next(context.Background())
if secondPage == nil {
    slog.Info("no more Pokémon")
	return
}
```

Alternatively, you can iterate through all Pokémon using the `All()` method on the returned paginator:
```go
client := pokesdk.NewClient()

// Fetch first page of Pokémon (use a better context in production code)
pages := client.Pokemon.List().All(context.Background())

for page := range pages {
	if page.Error != nil {
		fmt.Printf("Error fetching Pokémon list: %v", page.Error)
		os.Exit(1)
	}
	for _, pokemon := range page.Results {
		slog.Info("Pokémon: %s", pokemon.Name
		// you can get more information on the pokemon by fetching it by ref
		detailedPokemon, err := client.Pokemon.GetByRef(context.Background(), pokemon)
		if err != nil {
			fmt.Printf("Error fetching Pokémon details: %v", err)
			continue
		}
	}
}
```

### Fetching Pokémon by ID
```go
pokemon, err := client.Pokemon.GetByID(context.Background(), 25)
```
 
### Fetching Pokémon by Name
```go
pokemon, err := client.Pokemon.GetByName(context.Background(), "pikachu")
```
### Fetching Pokémon by Reference

>_NOTE: A reference is returned in the listing response, and isn't constructed by the user._

```go
ref := pokesdk.PokemonRef{
    Name: "pikachu",
    URL:  "https://pokeapi.co/api/v2/pokemon/25/",
}
pokemon, err := client.Pokemon.GetByRef(context.Background(), ref)

## Getting generation data

You can access generation data in a similar way to Pokémon data. Just use `client.Generation` instead of `client.Pokemon` in the above examples.

### Checking for errors
#### Not Found
If a resource is not found you can check for relevant not found error.

##### Pokémon Not Found

```go
_, err := client.Pokemon.GetByName(context.Background(), "nonexistent")
if err != nil {
	if errors.Is(err, pokesdk.PokemonNotFoundError) {
		fmt.Println("Pokémon not found")
	} else {
		fmt.Printf("Error fetching Pokémon: %v", err)
	}
	os.Exit(1)
}
```
#### Generation Not Found
```go
_, err := client.Generation.GetByID(context.Background(), 999)
if err != nil {
	if errors.Is(err, pokesdk.GenerationNotFoundError) {
		fmt.Println("Generation not found")
	} else {
		fmt.Printf("Error fetching generation: %v", err)
	}
	os.Exit(1)
}
```

## Running tests

You can use `make test-all` to run all tests, including unit and integration tests.

To run only unit tests, use `make test`, and to run only integration tests, use `make test-integration`.

## Example usage

You can see example usage of the SDK in [`cmd/example/main.go`](/cmd/example/main.go).

## Design decisions

1. I wanted to provide the ability for users to easily iterate through paginated results without having to manually handle pagination logic. The paginator allows users to decide whether they want to handle the concurrency themselves (using `Next()` method) or have the SDK handle it for them (using the `All()` method that returns a channel). I am concious of the fact that forcing concurrency on users is not ideal, and should be left up to the user to decide. I think this strikes a good balance.
2. I modelled the paginator as a generic type that can be used for any resource, allowing for code reuse and consistency across different resources. In the future adding new endpoints will be easier.
3. For integration tests I considered using something like WireMock and [test containers](https://golang.testcontainers.org) to spin it up as part of the test suite, but I decided to use a test HTTP server with simple stub responses instead. The main reason was to keep the test suite simple with as few external dependencies as possible.
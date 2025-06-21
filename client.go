package pokesdk

type Client struct {
	// Pokemon provides access to the Pokemon API endpoints
	Pokemon PokemonAPI
	// Generation provides access to the Generation API endpoints
	Generation GenerationAPI
}

func NewClient(opts ...Option) *Client {
	cfg := NewConfig(opts...)
	return &Client{
		Pokemon:    PokemonAPI{cfg: cfg},
		Generation: GenerationAPI{cfg: cfg},
	}
}

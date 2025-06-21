package pokesdk

type GenerationRef NamedAPIResource

type GenerationList struct {
	Count    int             `json:"count"`
	Next     *string         `json:"next"`
	Previous *string         `json:"previous"`
	Results  []GenerationRef `json:"results"`
}

func (g *GenerationList) GetNextURL() string {
	if g.Next == nil {
		return ""
	}
	return *g.Next
}

type Generation struct {
	ID             int                `json:"id"`
	Name           string             `json:"name"`
	Abilities      []NamedAPIResource `json:"abilities"`
	MainRegion     NamedAPIResource   `json:"main_region"`
	Moves          []NamedAPIResource `json:"moves"`
	Names          []LocalizedName    `json:"names"`
	PokemonSpecies []NamedAPIResource `json:"pokemon_species"`
	Types          []NamedAPIResource `json:"types"`
	VersionGroups  []NamedAPIResource `json:"version_groups"`
}

type LocalizedName struct {
	Language NamedAPIResource `json:"language"`
	Name     string           `json:"name"`
}

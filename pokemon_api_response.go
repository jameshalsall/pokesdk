package pokesdk

type PokemonRef NamedAPIResource

type PokemonList struct {
	Count    int          `json:"count"`
	Next     *string      `json:"next"`
	Previous *string      `json:"previous"`
	Results  []PokemonRef `json:"results"`
}

func (p *PokemonList) GetNextURL() string {
	if p.Next == nil {
		return ""
	}
	return *p.Next
}

type Pokemon struct {
	ID                     int                `json:"id"`
	Name                   string             `json:"name"`
	Abilities              []PokemonAbility   `json:"abilities"`
	BaseExperience         int                `json:"base_experience"`
	Cries                  PokemonCries       `json:"cries"`
	Forms                  []NamedAPIResource `json:"forms"`
	GameIndices            []GameIndex        `json:"game_indices"`
	Height                 int                `json:"height"`
	HeldItems              []HeldItems        `json:"held_items"`
	IsDefault              bool               `json:"is_default"`
	LocationAreaEncounters string             `json:"location_area_encounters"`
	Moves                  []PokemonMove      `json:"moves"`
	Order                  int                `json:"order"`
	PastAbilities          []PastAbilities    `json:"past_abilities"`
	PastTypes              []PastTypes        `json:"past_types"`
	Species                NamedAPIResource   `json:"species"`
	Sprites                PokemonSprites     `json:"sprites"`
	Stats                  []PokemonStat      `json:"stats"`
	Types                  []PokemonType      `json:"types"`
	Weight                 int                `json:"weight"`
}

type PokemonAbility struct {
	Ability  NamedAPIResource `json:"ability"`
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
}

type PokemonCries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

type GameIndex struct {
	GameIndex int              `json:"game_index"`
	Version   NamedAPIResource `json:"version"`
}

type PastTypes struct {
	Generation NamedAPIResource `json:"generation"`
	Types      []Types          `json:"types"`
}

type Types struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

type HeldItems struct {
	Item           NamedAPIResource `json:"item"`
	VersionDetails []VersionDetails `json:"version_details"`
}

type PokemonMove struct {
	Move                NamedAPIResource         `json:"move"`
	VersionGroupDetails []MoveVersionGroupDetail `json:"version_group_details"`
}

type MoveVersionGroupDetail struct {
	LevelLearnedAt  int              `json:"level_learned_at"`
	MoveLearnMethod NamedAPIResource `json:"move_learn_method"`
	Order           *int             `json:"order"`
	VersionGroup    NamedAPIResource `json:"version_group"`
}

type PastAbilities struct {
	Abilities  []PokemonAbility `json:"abilities"`
	Generation NamedAPIResource `json:"generation"`
}

type PokemonSprites struct {
	BackDefault      string   `json:"back_default"`
	BackFemale       *string  `json:"back_female"`
	BackShiny        string   `json:"back_shiny"`
	BackShinyFemale  *string  `json:"back_shiny_female"`
	FrontDefault     string   `json:"front_default"`
	FrontFemale      *string  `json:"front_female"`
	FrontShiny       string   `json:"front_shiny"`
	FrontShinyFemale *string  `json:"front_shiny_female"`
	Other            Others   `json:"other"`
	Versions         Versions `json:"versions"`
}

type Others map[string]map[string]string

type Versions map[string]map[string]map[string]any

type PokemonStat struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResource `json:"stat"`
}

type PokemonType struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

type VersionDetails struct {
	Rarity  int              `json:"rarity"`
	Version NamedAPIResource `json:"version"`
}

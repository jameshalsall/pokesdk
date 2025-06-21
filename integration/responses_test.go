package integration

import (
	"bytes"
	_ "embed"
	"text/template"
)

var (
	//go:embed testdata/pokemon_response.json
	pokemonResponse []byte
	//go:embed testdata/list_pokemon_response_page_1.json
	listPokemonResponsePage1 []byte
	//go:embed testdata/list_pokemon_response_page_2.json
	listPokemonResponsePage2 []byte

	//go:embed testdata/generation_response.json
	generationResponse []byte
	//go:embed testdata/list_generation_response_page_1.json
	listGenerationResponsePage1 []byte
	//go:embed testdata/list_generation_response_page_2.json
	listGenerationResponsePage2 []byte
)

func ResponseBytes(tpl []byte, baseURL string) []byte {
	data := map[string]string{
		"BaseURL": baseURL,
	}

	t, err := template.New("response").Parse(string(tpl))
	if err != nil {
		panic("failed to parse template: " + err.Error())
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		panic("failed to execute template: " + err.Error())
	}

	return buf.Bytes()
}

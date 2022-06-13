package util

import (
	"catching-pokemons/models"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParserSuccess(t *testing.T) {
	c := require.New(t)

	body, err := ioutil.ReadFile("samples/parsed_pokemon_response.json")
	c.NoError(err)

	var response models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &response)
	c.NoError(err)

	parsedPokemon, err := ParsePokemon(response)
	c.NoError(err)

	body, err = ioutil.ReadFile("samples/parsed_response.json")
	c.NoError(err)

	var expected models.Pokemon

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, parsedPokemon)
}

func TestParserPokemonTypeNotFound(t *testing.T) {
	c := require.New(t)

	body, err := ioutil.ReadFile("samples/parsed_pokemon_response.json")
	c.NoError(err)

	var response models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &response)
	c.NoError(err)

	response.PokemonType = []models.PokemonType{}

	_, err = ParsePokemon(response)
	c.NotNil(err)
	c.EqualError(ErrNotFoundPokemonType, err.Error())
}

func BenchmarkParser(b *testing.B) {
	c := require.New(b)

	body, err := ioutil.ReadFile("samples/parsed_pokemon_response.json")
	c.NoError(err)

	var response models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &response)
	c.NoError(err)

	body, err = ioutil.ReadFile("samples/parsed_response.json")
	c.NoError(err)

	var expected models.Pokemon

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	for n := 0; n < b.N; n++ {
		parsedPokemon, err := ParsePokemon(response)
		c.NoError(err)

		c.Equal(expected, parsedPokemon)
	}
}

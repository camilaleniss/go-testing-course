package controller

import (
	"catching-pokemons/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonFromPokeApiSuccess(t *testing.T) {
	c := require.New(t)

	pokemon, err := GetPokemonFromPokeApi("bulbasaur")
	c.NoError(err)

	body, err := ioutil.ReadFile("samples/parsed_pokemon_response.json")
	c.NoError(err)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, pokemon)
}

func TestGetPokemonFromPokeApiNotFound(t *testing.T) {
	c := require.New(t)

	_, err := GetPokemonFromPokeApi("aa")
	c.NotNil(err)
	c.EqualError(ErrPokemonNotFound, err.Error())
}

func TestGetPokemonFromPokeApiSuccessWithMocks(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "bulbasaur"

	body, err := ioutil.ReadFile("samples/parsed_pokemon_response.json")
	c.NoError(err)

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id),
		httpmock.NewStringResponder(200, string(body)))

	pokemon, err := GetPokemonFromPokeApi(id)
	c.NoError(err)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, pokemon)
}

func TestGetPokemonFromPokeApiNotFoundWithMocks(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "aaaaa"

	mockBody, err := ioutil.ReadFile("samples/pokemon_not_found_response.txt")
	c.NoError(err)

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id),
		httpmock.NewStringResponder(404, string(mockBody)))

	_, err = GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(ErrPokemonNotFound, err.Error())
}

func TestGetPokemonFromPokeApiInternalErrorWithMocks(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "aaaaa"

	mockBody, err := ioutil.ReadFile("samples/pokemon_not_found_response.txt")
	c.NoError(err)

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id),
		httpmock.NewStringResponder(500, string(mockBody)))

	_, err = GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(ErrPokeApiFailure, err.Error())
}

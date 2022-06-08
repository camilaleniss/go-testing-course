package controller

import (
	"catching-pokemons/models"
	"encoding/json"
	"io/ioutil"
	"testing"

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

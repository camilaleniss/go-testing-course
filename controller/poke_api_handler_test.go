package controller

import (
	"catching-pokemons/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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

func TestGetPokemon(t *testing.T) {
	c := require.New(t)

	r, _ := http.NewRequest("GET", "/pokemon/{id}", nil)
	w := httptest.NewRecorder()

	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": "bulbasaur",
	}

	r = mux.SetURLVars(r, vars)

	GetPokemon(w, r)

	expectedBodyResponse, err := ioutil.ReadFile("samples/parsed_response.json")
	c.NoError(err)

	var expectedPokemon models.Pokemon

	err = json.Unmarshal([]byte(expectedBodyResponse), &expectedPokemon)
	c.NoError(err)

	var actualPokemon models.Pokemon

	err = json.Unmarshal(w.Body.Bytes(), &actualPokemon)
	c.NoError(err)

	c.Equal(http.StatusOK, w.Code)
	c.Equal(expectedPokemon, actualPokemon)
}

package poke

import (
	"testing"

	"gotest.tools/assert"
)

func TestInitializeMapTypeCheck(t *testing.T) {
	var listType []string
	initializeMapTypeCheck()
	listType = append(listType, "error")
	accepted, _, _, _ := checkAcceptedType(listType)
	assert.Equal(t, accepted, false)
	listType = append(listType, "fire")
	accepted, _, _, _ = checkAcceptedType(listType)
	assert.Equal(t, accepted, true)
	listType = append(listType, "flying")
	_, acceptedExtra, _, _ := checkAcceptedType(listType)
	assert.Equal(t, acceptedExtra, true)
}

func TestRetrievePokemonType(t *testing.T) {
	pokemonImpl := NewPokemonAPIimpl()
	accepted, _, _, _, _ := pokemonImpl.RetrievePokemonType("bulbasaur")
	assert.Equal(t, accepted, true)
	accepted, _, _, _, _ = pokemonImpl.RetrievePokemonType("pikachu")
	assert.Equal(t, accepted, false)
	_, _, _, _, err := pokemonImpl.RetrievePokemonType("ERROR")
	assert.Equal(t, err != nil, true)
}

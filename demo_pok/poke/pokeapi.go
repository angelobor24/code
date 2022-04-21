package poke

import (
	"fmt"

	"github.com/mtslzr/pokeapi-go"
)

var (
	mapTypeBase  = make(map[string]bool)
	mapTypeExtra = make(map[string]bool)
	// list for accepted category pokemon to start an insurance
	listAcceptedTypes = []string{"fire", "water", "grass"}
	// list of category pokemon with extra price
	listAcceptedExtra = []string{"flying"}
)

type PokemonAPI interface {
	RetrievePokemonType(name string) (bool, bool, string, string, error)
}

type PokemonAPIimpl struct {
}

func NewPokemonAPIimpl() PokemonAPI {
	pokeImpl := PokemonAPIimpl{}
	initializeMapTypeCheck()
	return &pokeImpl
}

// Check if the provided pokemon name exist, and retrieve the related category
// using the pokapi client. If the pokemon exist, return if the category is valid,
// considering the listAcceptedTypes and the listAcceptedExtra. If yes, also the category pokemon type
// are returned.

func (pokemonImpl *PokemonAPIimpl) RetrievePokemonType(name string) (bool, bool, string, string, error) {
	var listPokemonType []string
	var typeBase, typeExtra string
	types, err := pokeapi.Pokemon(name)
	if err != nil {
		fmt.Println(err)
		return false, false, typeBase, typeExtra, err
	}
	typesPokemon := types.Types
	for _, s := range typesPokemon {
		listPokemonType = append(listPokemonType, s.Type.Name)
	}
	existBaseType, existExtraType, typeBase, typeExtra := checkAcceptedType(listPokemonType)
	return existBaseType, existExtraType, typeBase, typeExtra, nil
}

// check if the retrieved pokemon type list is valid considering the specific
// list of accepted type.
func checkAcceptedType(listTypes []string) (bool, bool, string, string) {
	acceptedTypeBase := false
	acceptedTypeExtra := false
	var typeBase, typeExtra string
	for _, typeRetrieved := range listTypes {
		if !acceptedTypeBase {
			if mapTypeBase[typeRetrieved] {
				typeBase = typeRetrieved
				acceptedTypeBase = true
			}
		}
		if !acceptedTypeExtra {
			if mapTypeExtra[typeRetrieved] {
				typeExtra = typeRetrieved
				acceptedTypeExtra = true
			}
		}

	}
	return acceptedTypeBase, acceptedTypeExtra, typeBase, typeExtra
}

// initialize the map used to check if a category is accepted
func initializeMapTypeCheck() {
	for _, acceptedType := range listAcceptedTypes {
		mapTypeBase[acceptedType] = true
	}
	for _, acceptedType := range listAcceptedExtra {
		mapTypeExtra[acceptedType] = true
	}
}

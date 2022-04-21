package main

import (
	"demo/db"
	"demo/payment"
	"demo/poke"
	"demo/server"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting service..")

	//structure initialization
	pokeApi := poke.NewPokemonAPIimpl()
	dbImpl := db.NewStorageImpl("sqlite3", "pokemonInsurance", "?_foreign_keys=true")
	if !dbImpl.GetStatus() {
		fmt.Println("DB UNAVAILABLE. The service isn't available to start..")
		os.Exit(1)
	}
	dbImpl.InitializeTables(db.ListTables[:])

	// Insert price for the approved insurance category
	dbImpl.InsertBaseQuote("fire", 9.5)
	dbImpl.InsertBaseQuote("grass", 9.5)
	dbImpl.InsertBaseQuote("water", 9.5)
	dbImpl.InsertExtraQuote("flying", 0.5)

	// Start server
	server := server.NewServerImpl(server.NewServiceImpl(pokeApi, payment.NewPaymentImpl()), dbImpl)
	server.StartServer()
}

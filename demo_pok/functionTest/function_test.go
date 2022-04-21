package functionTest

import (
	"bytes"
	"demo/db"
	"demo/payment"
	"demo/poke"
	"demo/server"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"gotest.tools/assert"
)

/*
Test specification:
1) Register a new trainer: success operation
2) Register the same trainer: error during operation, because the trainer is already registered
3) Performs a quote for a grass pokemon: success operation, the total price is 9.5
4) Performs the same quote: error during operation, because the quote is already done for the pokemon
5) Performs a quote for a fire & flying pokemon: success operation, the total price is 10
6) Pay the quote related to grass pokemon: success operation
7) Pay the quote for a grass pokemon for which quote is not prepared: failed operation
8) Performs a quote for a fire & flying pokemon: success operation, the total price is 10
9) Disable the payment system. Pay the quote for the quote performed to point 8): failed operation, payment rejected.
*/

func TestFT1(t *testing.T) {
	var totalPrice float32 = 9.5
	pokeApi := poke.NewPokemonAPIimpl()
	dbImpl := db.NewStorageImpl("sqlite3", "pokemonInsurance", "?_foreign_keys=true")
	dbImpl.InitializeTables(db.ListTables[:])
	serverImpl := server.NewServerImpl(server.NewServiceImpl(pokeApi, payment.NewPaymentImpl()), dbImpl)
	go serverImpl.StartServer()
	time.Sleep(2 * time.Second)
	trainer := server.TrainerInfo{}
	trainer.Idtrainer = 24
	trainer.Name = "testq"
	trainer.Surname = "testq"
	quote := server.Quote{}
	quote.Pokemon = "bulbasaur"
	postBody, _ := json.Marshal(trainer)
	responseBody := bytes.NewBuffer(postBody)
	resp, _ := http.Post("http://127.0.0.1:8080/trainer", "application/json", responseBody)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	responseBody = bytes.NewBuffer(postBody)
	resp, _ = http.Post("http://127.0.0.1:8080/trainer", "application/json", responseBody)
	assert.Equal(t, resp.StatusCode, http.StatusConflict)
	dbImpl.InsertBaseQuote("fire", 9.5)
	dbImpl.InsertBaseQuote("grass", 9.5)
	dbImpl.InsertBaseQuote("water", 9.5)
	dbImpl.InsertExtraQuote("flying", 0.5)

	postBody, _ = json.Marshal(quote)
	responseBody = bytes.NewBuffer(postBody)
	newQuote := server.Quote{}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/quote", responseBody)
	req.Header.Add("mt", "25")
	resp, _ = client.Do(req)
	json.NewDecoder(resp.Body).Decode(&newQuote)
	assert.Equal(t, newQuote.Pokemon, "bulbasaur")
	assert.Equal(t, newQuote.Price, totalPrice)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)

	postBody, _ = json.Marshal(quote)
	responseBody = bytes.NewBuffer(postBody)
	req, _ = http.NewRequest("POST", "http://127.0.0.1:8080/quote", responseBody)
	req.Header.Add("mt", "25")
	resp, _ = client.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusConflict)

	quote.Pokemon = "charizard"
	postBody, _ = json.Marshal(quote)
	responseBody = bytes.NewBuffer(postBody)

	req, _ = http.NewRequest("POST", "http://127.0.0.1:8080/quote", responseBody)
	req.Header.Add("mt", "25")
	resp, _ = client.Do(req)
	json.NewDecoder(resp.Body).Decode(&newQuote)
	totalPrice = 10.0
	assert.Equal(t, newQuote.Price, totalPrice)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)

	payedQuote := server.PayedQuote{}
	payedQuote.Pokemon = "charizard"
	postBody, _ = json.Marshal(payedQuote)
	responseBody = bytes.NewBuffer(postBody)

	req, _ = http.NewRequest("POST", "http://127.0.0.1:8080/payedQuote", responseBody)
	req.Header.Add("mt", "25")
	resp, _ = client.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)

	postBody, _ = json.Marshal(payedQuote)
	responseBody = bytes.NewBuffer(postBody)

	req, _ = http.NewRequest("POST", "http://127.0.0.1:8080/payedQuote", responseBody)
	req.Header.Add("mt", "25")
	resp, _ = client.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	os.Setenv("PAYMENT_BEHAVIOUR", "rejected")

	postBody, _ = json.Marshal(payedQuote)
	responseBody = bytes.NewBuffer(postBody)

	req, _ = http.NewRequest("POST", "http://127.0.0.1:8080/quote", responseBody)
	req.Header.Add("mt", "25")
	resp, _ = client.Do(req)
	json.NewDecoder(resp.Body).Decode(&newQuote)
	totalPrice = 10.0
	assert.Equal(t, newQuote.Price, totalPrice)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)

	payedQuote.Pokemon = "charizard"
	postBody, _ = json.Marshal(payedQuote)
	responseBody = bytes.NewBuffer(postBody)

	req, _ = http.NewRequest("POST", "http://127.0.0.1:8080/payedQuote", responseBody)
	req.Header.Add("mt", "25")
	resp, _ = client.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM quote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM extraquote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM payedquote")
	os.Unsetenv("PAYMENT_BEHAVIOUR")
	defer resp.Body.Close()
}

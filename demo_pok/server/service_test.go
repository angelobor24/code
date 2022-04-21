package server

import (
	"demo/db"
	"demo/handlerMessage"
	"demo/payment"
	"demo/poke"
	"testing"

	"gotest.tools/assert"
)

type Test struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Idtrainer int    `json:"idtrainer"`
}

func TestAddTrainer(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var trainer TrainerInfo
	trainer.Name = "test"
	trainer.Surname = "test"
	trainer.Idtrainer = 1
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	_, err := serviceImpl.AddTrainer(trainer, dbImpl)
	assert.Equal(t, err, nil)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
}

func TestAddTrainerInternalError(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var trainer TrainerInfo
	trainer.Name = "test"
	trainer.Surname = "test"
	trainer.Idtrainer = 1
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	_, err := serviceImpl.AddTrainer(trainer, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrInternalError)
}

func TestAddTrainerInputError(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var trainer TrainerInfo
	trainer.Name = "test"
	trainer.Idtrainer = 1
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	_, err := serviceImpl.AddTrainer(trainer, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrInputParameters)
}

func TestAddTrainerResourceAlreadyExist(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var trainer TrainerInfo
	trainer.Name = "test"
	trainer.Surname = "test"
	trainer.Idtrainer = 1
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	_, err := serviceImpl.AddTrainer(trainer, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrResourceAlreadyExist)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
}

func TestAddQuote(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "test"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	_, err := serviceImpl.AddQuote(quote, 2, dbImpl)
	assert.Equal(t, err != nil, true)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
}

func TestAddQuoteMtNotExist(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "test"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	_, err := serviceImpl.AddQuote(quote, 2, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrSecurityToken)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
}

func TestAddQuoteErrorInput(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	_, err := serviceImpl.AddQuote(quote, 2, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrInputParameters)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
}

func TestAddQuoteValidPoke(t *testing.T) {
	var price float32 = 9.5
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "squirtle"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('water',9.5)")
	quote, err := serviceImpl.AddQuote(quote, 2, dbImpl)
	assert.Equal(t, err, nil)
	assert.Equal(t, quote.Price, price)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
}

func TestAddQuoteValidPokeExtraQuote(t *testing.T) {
	var price float32 = 10.0
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "charizard"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM extraquote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('fire',9.5)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO extraquote VALUES ('flying',0.5)")
	quoteResult, err := serviceImpl.AddQuote(quote, 2, dbImpl)
	assert.Equal(t, err, nil)
	assert.Equal(t, quoteResult.Price, price)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE extraquote")
}

func TestAddQuoteValidPokeQuoteAlreadyExist(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "squirtle"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM quote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('water',9.5)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO quote VALUES (1,'squirtle',2)")
	_, err := serviceImpl.AddQuote(quote, 2, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrResourceAlreadyExist)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
}

func TestAddQuoteInvalidCategory(t *testing.T) {
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "pikachu"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('water',9.5)")
	_, err := serviceImpl.AddQuote(quote, 2, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrPokemonCategory)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
}

func TestAddPayedQuote(t *testing.T) {
	payedQuote := PayedQuote{}
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "squirtle"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, db.CreatePayedQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM quote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM payedquote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('water',9.5)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO quote VALUES (1,'bulbasaur',2)")
	payedQuote.Pokemon = "bulbasaur"
	serviceImpl.AddQuote(quote, 2, dbImpl)
	_, err := serviceImpl.AddInsurance(payedQuote, 2, dbImpl)
	assert.Equal(t, err, nil)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE payedquote")
}

func TestAddPayedQuoteInvalidMt(t *testing.T) {
	payedQuote := PayedQuote{}
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "squirtle"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, db.CreatePayedQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM quote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM payedquote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('water',9.5)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO quote VALUES (1,'bulbasaur',2)")
	payedQuote.Pokemon = "bulbasaur"
	serviceImpl.AddQuote(quote, 2, dbImpl)
	_, err := serviceImpl.AddInsurance(payedQuote, 20, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrSecurityToken)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE payedquote")
}

func TestAddPayedQuoteNoteFound(t *testing.T) {
	payedQuote := PayedQuote{}
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "squirtle"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, db.CreatePayedQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM quote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM payedquote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('water',9.5)")
	payedQuote.Pokemon = "bulbasaur"
	serviceImpl.AddQuote(quote, 2, dbImpl)
	_, err := serviceImpl.AddInsurance(payedQuote, 2, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrQuoteNotFound)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE payedquote")
}

func TestAddPayedQuoteErrorInput(t *testing.T) {
	payedQuote := PayedQuote{}
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(true)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "squirtle"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, db.CreatePayedQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM quote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM payedquote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('water',9.5)")
	serviceImpl.AddQuote(quote, 2, dbImpl)
	_, err := serviceImpl.AddInsurance(payedQuote, 2, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrInputParameters)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE payedquote")
}

func TestAddPayedQuoteRejectedPayment(t *testing.T) {
	payedQuote := PayedQuote{}
	pokeApi := poke.NewPokemonAPIimpl()
	paymentImpl := payment.MockNewPaymentImpl(false)
	serviceImpl := NewServiceImpl(pokeApi, paymentImpl)
	var quote Quote
	quote.Pokemon = "squirtle"
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	db.ActionOnTable(db.DbInstance, db.CreateTrainer)
	db.ActionOnTable(db.DbInstance, db.CreateQuote)
	db.ActionOnTable(db.DbInstance, db.CreateBaseQuote)
	db.ActionOnTable(db.DbInstance, db.CreateExtraQuote)
	db.ActionOnTable(db.DbInstance, db.CreatePayedQuote)
	db.ActionOnTable(db.DbInstance, "DELETE FROM trainer")
	db.ActionOnTable(db.DbInstance, "DELETE FROM quote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM basequote")
	db.ActionOnTable(db.DbInstance, "DELETE FROM payedquote")
	db.ActionOnTable(db.DbInstance, "INSERT INTO trainer VALUES (1,'test','test',2)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO basequote VALUES ('water',9.5)")
	db.ActionOnTable(db.DbInstance, "INSERT INTO quote VALUES (1,'bulbasaur',2)")
	payedQuote.Pokemon = "bulbasaur"
	serviceImpl.AddQuote(quote, 2, dbImpl)
	_, err := serviceImpl.AddInsurance(payedQuote, 2, dbImpl)
	assert.Equal(t, err, handlerMessage.ErrPaymentSystem)
	db.ActionOnTable(db.DbInstance, "DROP TABLE trainer")
	db.ActionOnTable(db.DbInstance, "DROP TABLE quote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE basequote")
	db.ActionOnTable(db.DbInstance, "DROP TABLE payedquote")
}

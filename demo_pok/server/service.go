package server

import (
	"demo/db"
	"demo/handlerMessage"
	"demo/payment"
	"demo/poke"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type TrainerInfo struct {
	Name      string `json:"name" validate:"required,alphaunicode"`
	Surname   string `json:"surname" validate:"required,alphaunicode"`
	Idtrainer int    `json:"idtrainer" validate:"required"`
	Mt        int    `json:"mt" `
}

type Quote struct {
	Pokemon string  `json:"pokemon" validate:"required"`
	Price   float32 `json:"price" `
	Id      int     `json:"id" `
}

type PayedQuote struct {
	Pokemon       string  `json:"pokemon" validate:"required"`
	Price         float32 `json:"price" `
	IdTrainer     int     `json:"idTrainer" `
	IdTransaction int64   `json:"idTransaction" `
	Timestamp     string  `json:"timestamp" `
}
type Service interface {
	AddTrainer(TrainerInfo, db.Storage) (TrainerInfo, error)
	AddQuote(Quote, int, db.Storage) (Quote, error)
	AddInsurance(payedQuote PayedQuote, mt int, storage db.Storage) (PayedQuote, error)
}

type ServiceImpl struct {
	pokeApi       poke.PokemonAPI
	paymentSystem payment.Payment
}

// service implementer of the Service interface. This structure implements all the logicto handle the
// received server request
func NewServiceImpl(pokeApi poke.PokemonAPI, paymentSystem payment.Payment) Service {
	serviceImpl := ServiceImpl{pokeApi: pokeApi, paymentSystem: paymentSystem}
	return &serviceImpl
}

// service to add a new trainer into the system
func (serviceImpl *ServiceImpl) AddTrainer(trainer TrainerInfo, storage db.Storage) (TrainerInfo, error) {
	validate := validator.New()
	// validate the input body
	err := validate.Struct(trainer)
	if err != nil {
		fmt.Println(err)
		return trainer, handlerMessage.ErrInputParameters
	}
	results, err := storage.CheckTrainerExist(trainer.Idtrainer)
	if err != nil {
		fmt.Println(err)
		return trainer, handlerMessage.ErrInternalError
	}
	defer results.Close()

	// check if the trainer is already registered
	if !results.Next() {
		// associate an unique mt to the trainer
		trainer.Mt = createRandomUniqueMt(trainer.Idtrainer)
		err = storage.InsertTrainer(trainer.Idtrainer, trainer.Name, trainer.Surname, trainer.Mt)
		if err != nil {
			fmt.Println(err)
			return trainer, handlerMessage.ErrInternalError
		}
	} else {
		// trainer already registered
		return trainer, handlerMessage.ErrResourceAlreadyExist
	}
	return trainer, nil
}
func (serviceImpl *ServiceImpl) AddQuote(quote Quote, mt int, storage db.Storage) (Quote, error) {
	idTrainer, exist := checkConsistencyMt(storage, mt)
	// check if the received mt is valid
	if !exist {
		fmt.Println("Received invalid mt")
		return quote, handlerMessage.ErrSecurityToken
	}
	// validate the input body request
	validate := validator.New()
	err := validate.Struct(quote)
	if err != nil {
		fmt.Println(err)
		return quote, handlerMessage.ErrInputParameters
	}
	quote.Id = idTrainer
	// calculate a quote for the pokemon
	totalPrice, err := serviceImpl.retrievePricePokemon(storage, quote.Pokemon)
	quote.Price = totalPrice
	if err != nil {
		// error during interaction with poke api server
		fmt.Println(err)
		return quote, err
	}

	// add new quote
	err = storage.InsertQuote(quote.Id, quote.Pokemon, quote.Price)
	if err != nil {
		fmt.Println(err)
		return quote, handlerMessage.ErrResourceAlreadyExist
	}
	return quote, nil
}
func (serviceImpl *ServiceImpl) AddInsurance(payedQuote PayedQuote, mt int, storage db.Storage) (PayedQuote, error) {
	idTrainer, exist := checkConsistencyMt(storage, mt)
	// check if the received mt is valid
	if !exist {
		fmt.Println("Received invalid mt")
		return payedQuote, handlerMessage.ErrSecurityToken
	}
	//validate the input body request
	validate := validator.New()
	err := validate.Struct(payedQuote)
	if err != nil {
		fmt.Println(err)
		return payedQuote, handlerMessage.ErrInputParameters
	}

	payedQuote.IdTrainer = idTrainer

	// retrieve the related to quote before to pay
	results, err := storage.GetFromQuoteFiltered(idTrainer, payedQuote.Pokemon)
	if err != nil {
		fmt.Println(err)
		return payedQuote, handlerMessage.ErrInternalError
	}

	if !results.Next() {
		// doesn't exist a quote to pay
		return payedQuote, handlerMessage.ErrQuoteNotFound
	} else {
		results.Scan(&payedQuote.Price)
	}
	results.Close()

	// Pay the quote
	if serviceImpl.paymentSystem.Pay() {
		transaction, _ := db.ReturnTransactionContext()
		// remove the quote
		err = storage.DeleteFromQuote(idTrainer, payedQuote.Pokemon)
		if err != nil {
			fmt.Println(err)
			return payedQuote, handlerMessage.ErrInternalError
		}
		// generate an unique Id for the transaction
		uniqueId := returnUniqueRandomId()
		err = storage.InsertPayedQuote(uniqueId, idTrainer, payedQuote.Pokemon, payedQuote.Price)
		if err != nil {
			fmt.Println(err)
			transaction.Rollback()
			return payedQuote, handlerMessage.ErrInternalError
		}
		transaction.Commit()
		t := time.Now()
		payedQuote.IdTransaction = uniqueId
		payedQuote.Timestamp = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		//

		return payedQuote, nil
	}

	// error during payment
	return payedQuote, handlerMessage.ErrPaymentSystem
}

func createRandomUniqueMt(variableNumber int) int {
	return variableNumber + 1
}

func checkConsistencyMt(storage db.Storage, mt int) (int, bool) {
	idTrainer, err := storage.CheckTrainerMt(mt)
	if err != nil {
		fmt.Println(err)
		return idTrainer, false
	}
	return idTrainer, true
}

func (serviceImpl *ServiceImpl) retrievePricePokemon(storage db.Storage, pokemon string) (float32, error) {
	var totalPrice, partialPrice float32
	isvalidCategoryBase, isvalidCategoryExtra, categoryBase, categoryExtra, err := serviceImpl.pokeApi.RetrievePokemonType(pokemon)
	if err != nil {
		fmt.Println(err)
		return totalPrice, handlerMessage.ErrServicePokemon
	}
	if !isvalidCategoryBase {
		return totalPrice, handlerMessage.ErrPokemonCategory
	}
	results, err := storage.GetFromBaseQuote(categoryBase)
	defer results.Close()

	if results.Next() {
		results.Scan(&partialPrice)
	}
	totalPrice = totalPrice + partialPrice
	if isvalidCategoryExtra {
		results, err = storage.GetFromExtraQuote(categoryExtra)
		defer results.Close()
		if results.Next() {
			results.Scan(&partialPrice)
		}
		totalPrice = totalPrice + partialPrice
	}
	return totalPrice, err
}

func returnUniqueRandomId() int64 {
	return time.Now().UnixNano() / (1 << 22)
}

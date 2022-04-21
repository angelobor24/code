package server

import (
	"demo/db"
)

type MockServiceImpl struct {
	trainerInfo TrainerInfo
	quote       Quote
	payedQuote  PayedQuote
	err         error
}

func NewMockServiceImpl(trainerInfo TrainerInfo, err error, quote Quote, payedQuote PayedQuote) Service {
	serviceMockImpl := MockServiceImpl{trainerInfo: trainerInfo, err: err, quote: quote}
	return &serviceMockImpl
}

func (mockServiceImpl *MockServiceImpl) AddTrainer(trainer TrainerInfo, storage db.Storage) (TrainerInfo, error) {
	return mockServiceImpl.trainerInfo, mockServiceImpl.err
}

func (mockServiceImpl *MockServiceImpl) AddQuote(quote Quote, mt int, storage db.Storage) (Quote, error) {
	return mockServiceImpl.quote, mockServiceImpl.err
}

func (mockServiceImpl *MockServiceImpl) AddInsurance(payedQuote PayedQuote, mt int, storage db.Storage) (PayedQuote, error) {
	return mockServiceImpl.payedQuote, mockServiceImpl.err
}

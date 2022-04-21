package server

import (
	"bytes"
	"context"
	"demo/db"
	"demo/handlerMessage"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func NewServerImplTest(service Service, storage db.Storage) ServerImpl {
	serverImpl := ServerImpl{service: service, storage: storage}
	return serverImpl
}

func TestQuote(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ("/quote"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	req.Header.Add("mt", "test")
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, nil, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.quote(w, req)
	assert.Equal(t, w.Code, http.StatusCreated)
}

func TestQuoteError(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ("/quote"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	req.Header.Add("mt", "test")
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, handlerMessage.ErrInternalError, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.quote(w, req)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestQuoteMethodNotAllowed(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, ("/quote"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	req.Header.Add("mt", "test")
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, handlerMessage.ErrInternalError, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.quote(w, req)
	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func TestQuotePathNotFound(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, ("/quoteTest"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	req.Header.Add("mt", "test")
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, handlerMessage.ErrInternalError, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.quote(w, req)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestTrainer(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ("/trainer"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, nil, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.trainer(w, req)
	assert.Equal(t, w.Code, http.StatusCreated)
}

func TestTrainerError(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/trainer", bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, handlerMessage.ErrInternalError, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.trainer(w, req)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestTrainerMethodNotAllowed(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, ("/trainer"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, handlerMessage.ErrInternalError, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.trainer(w, req)
	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func TestTrainerPathNotFound(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprint("/trainerTest"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, handlerMessage.ErrInternalError, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.trainer(w, req)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestPayedQuote(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ("/payedQuote"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, nil, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.payedquote(w, req)
	assert.Equal(t, w.Code, http.StatusCreated)
}
func TestPayedQuoteError(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ("/payedQuote"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, handlerMessage.ErrInternalError, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.payedquote(w, req)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestPayedQuoteMethodNotAllowed(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, ("/payedQuote"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, nil, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.payedquote(w, req)
	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func TestPayedQuotePathNotFound(t *testing.T) {
	dbImpl := db.NewStorageImpl("sqlite3", "dummy", "?_foreign_keys=true")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, ("/payedQuoteTest"), bytes.NewBuffer([]byte(`{"name":"test","surname":"test","idtrainer":45}`)))
	assert.NilError(t, err)
	w := httptest.NewRecorder()
	serviceMock := NewMockServiceImpl(TrainerInfo{Name: "test", Surname: "test", Idtrainer: 24}, nil, Quote{Pokemon: "test", Price: 0.0, Id: 1}, PayedQuote{Pokemon: "test", Price: 0.0, IdTrainer: 0, IdTransaction: 0, Timestamp: "test"})
	server := NewServerImplTest(serviceMock, dbImpl)
	server.payedquote(w, req)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

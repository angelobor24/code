package server

import (
	"demo/db"
	"demo/handlerMessage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server interface {
	StartServer()
}

type ServerImpl struct {
	service Service
	storage db.Storage
}

func NewServerImpl(service Service, storage db.Storage) Server {
	serverImpl := ServerImpl{service: service, storage: storage}
	return &serverImpl
}

func (serverImpl *ServerImpl) StartServer() {

	http.HandleFunc("/trainer", serverImpl.trainer)       // POST
	http.HandleFunc("/quote", serverImpl.quote)           // POST
	http.HandleFunc("/payedQuote", serverImpl.payedquote) // POST

	fmt.Printf("Starting insurance pokemon service at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// handler for POST request to /trainer endpoint
func (serverImpl *ServerImpl) trainer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/trainer" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		trainer := TrainerInfo{}
		json.NewDecoder(r.Body).Decode(&trainer)
		resourceCreated, err := serverImpl.service.AddTrainer(trainer, serverImpl.storage)
		if err != nil {
			statusCode, errorMessage := handlerMessage.ToStatusCodeMessage(err)
			http.Error(w, errorMessage, statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resourceCreated)
		return
	}

	http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
}

// handler for POST request to /quote endpoint
func (serverImpl *ServerImpl) quote(w http.ResponseWriter, r *http.Request) {
	// read mt from header request
	check := r.Header.Get("mt")
	if r.URL.Path != "/quote" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		quote := Quote{}
		mt, _ := strconv.Atoi(check)
		json.NewDecoder(r.Body).Decode(&quote)
		newQuote, err := serverImpl.service.AddQuote(quote, mt, serverImpl.storage)
		if err != nil {
			statusCode, errorMessage := handlerMessage.ToStatusCodeMessage(err)
			http.Error(w, errorMessage, statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newQuote)
		return
	}

	http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
}

// handler for POST request to /payedquote endpoint
func (serverImpl *ServerImpl) payedquote(w http.ResponseWriter, r *http.Request) {
	// read mt from header request
	check := r.Header.Get("mt")
	if r.URL.Path != "/payedQuote" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		mt, _ := strconv.Atoi(check)
		payedQuote := PayedQuote{}
		json.NewDecoder(r.Body).Decode(&payedQuote)
		newPayedQuote, err := serverImpl.service.AddInsurance(payedQuote, mt, serverImpl.storage)
		if err != nil {
			statusCode, errorMessage := handlerMessage.ToStatusCodeMessage(err)
			http.Error(w, errorMessage, statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPayedQuote)
		return
	}

	http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
}

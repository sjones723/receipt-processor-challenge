package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

// encode a given (v)alue as JSON ~ json.NewEncoder
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type apiError struct { Error string }

//  allows the use of ordinary functions as HTTP handlers 
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w, r); err != nil {
			//handle err
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

// controller function to start server 
func (s *APIServer) Run() {
	// mux router
	router := mux.NewRouter()

	router.HandleFunc("/receipts/process", makeHTTPHandleFunc(s.handleReceipt))

	log.Println("JSON API server running on port: ", s.listenAddr)
	
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleReceipt(w http.ResponseWriter, r *http.Request) error {

	// mux does not handle http type, must be handled in function
	if r.Method == "GET" {
		return s.handleGetPoints(w,r)
	}

	if r.Method == "POST" {
		return s.handleCreateUUID(w,r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateUUID(w http.ResponseWriter, r *http.Request) error {
	
	// return { "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
	return nil
}

func (s *APIServer) handleGetPoints(w http.ResponseWriter, r *http.Request) error {
	
	// return { "points": 32 }
	return nil
}









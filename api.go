package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
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

// controller, function to start server 
func (s *APIServer) Run() {
	// mux router
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(s.handleReturnAllReceipts))

	router.HandleFunc("/receipts/process", makeHTTPHandleFunc(s.handleAssignUUID))

	router.HandleFunc("/receipts/{id}/points", makeHTTPHandleFunc(s.handleCalculatePoints))

	log.Println("JSON API server running on port: ", s.listenAddr)
	
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleReturnAllReceipts(w http.ResponseWriter, r *http.Request) error {
	// Return a list of all receipts
	var receipts []Receipt
	for _, receipt := range receiptStore {
		receipts = append(receipts, receipt)
	}

	// Return the response as JSON
	return WriteJSON(w, http.StatusOK, receipts)
}

// Implement the /receipts/process endpoint
// @Summary Submits a receipt for processing
// @Description Submits a receipt for processing
// @Accept json
// @Produce json
// @Param req body Receipt true "Receipt object to process"
// @Success 200 {object} map[string]string
// @Failure 400 {object} apiError
// @Router /receipts/process [post]
func (s *APIServer) handleAssignUUID(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	id := uuid.New().String()
    response := map[string]string{"id": id}

	var newReceipt Receipt
		err := json.NewDecoder(r.Body).Decode(&newReceipt)
		if err != nil {
			return err
		}

		// Assign a new UUID to the Receipt
		newReceipt.ID = uuid.New().String()

		// Store the receipt locally (in-memory map)
		receiptStore[newReceipt.ID] = newReceipt


    return WriteJSON(w, http.StatusOK, response)


}

// Implement the /receipts/{id}/points endpoint
// @Summary Returns the points awarded for the receipt
// @Description Returns the points awarded for the receipt
// @Accept json
// @Produce json
// @Param id path string true "ID of the receipt"
// @Success 200 {object} map[string]int
// @Failure 404 {object} apiError
// @Router /receipts/{id}/points [get]
func (s *APIServer) handleCalculatePoints(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	// Parse the JSON receipt from the request body
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		return err
	}

	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return err
	}
	if total == float64(int(total)) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	numItems := len(receipt.Items)
	points += (numItems / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return err
			}
			itemPoints := int(math.Ceil(price * 0.2))
			points += itemPoints
		}
	}
	

	// Rule 6: 6 points if the day in the purchase date is odd
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return err
	}
	if purchaseDate.Day()%2 == 1 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return err
	}
	if purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) && purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
		points += 10
	}

	// Create a response JSON object with the calculated points
	response := map[string]int{"points": points}

	// Return the response as JSON
	return WriteJSON(w, http.StatusOK, response)
}










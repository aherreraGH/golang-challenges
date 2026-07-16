package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"localpractice7.com/challenges/internal/types"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Inventory Management", r.URL.Path[1:])
}

func inventoryAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(svc.Get()); err != nil {
		http.Error(w, "failed to encode inventory", http.StatusInternalServerError)
		return
	}
}

func inventoryFindHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	book, err := svc.Find(id)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "book not found in inventory",
		})
	}
	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "failed to encode inventory", http.StatusInternalServerError)
		return
	}
}

func inventoryCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data types.PostBookRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "bad data submitted",
		})
	}

	book, err := svc.CheckOut(data.Code, data.Who)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "book not found in inventory",
		})
		return
	}
	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "failed to encode inventory", http.StatusInternalServerError)
		return
	}

}

func inventoryStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(svc.GetStatistics()); err != nil {
		http.Error(
			w,
			"failed to encode statistics",
			http.StatusInternalServerError,
		)
		return
	}

}

func inventoryCheckInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data types.PostBookRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "bad data submitted",
		})
	}

	book, err := svc.CheckIn(data.Code)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "book not found in inventory",
		})
		return
	}
	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "failed to encode inventory", http.StatusInternalServerError)
		return
	}

}

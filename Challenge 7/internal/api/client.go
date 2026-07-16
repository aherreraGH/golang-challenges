package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"localpractice7.com/challenges/internal/inventory"
)

var svc inventory.Inventory

func main() {
	ctx := context.Background()

	svc.Load(ctx, "./data/inventory.json")

	// handle each path in the URL/URI
	http.HandleFunc("/", handler)
	http.HandleFunc("/inventory", inventoryAllHandler)
	http.HandleFunc("/inventory/{id}", inventoryFindHandler)
	http.HandleFunc("/checkout", inventoryCheckoutHandler)
	http.HandleFunc("/checkin", inventoryCheckInHandler)
	http.HandleFunc("/stats", inventoryStatsHandler)
	// show some output, otherwise it will appear as if nothing has happened
	fmt.Println("Staring server at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

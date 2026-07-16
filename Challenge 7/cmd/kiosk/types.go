package main

type Config struct {
	InventoryURI  string `json:"inventory-uri"`
	InventoryPort uint32 `json:"inventory-port"`
}

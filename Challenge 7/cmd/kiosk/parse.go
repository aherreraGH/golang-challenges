package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"localpractice7.com/challenges/internal/barcode"
	"localpractice7.com/challenges/internal/inventory"
	"localpractice7.com/challenges/internal/types"
)

func ListInventory(ctx context.Context) []inventory.Book {
	resp, err := http.Get(fmt.Sprintf("%s:%d/inventory", appConfig.Config.InventoryURI, appConfig.Config.InventoryPort))
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	var books []inventory.Book

	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	return books

}

func GetInventoryStatistics(ctx context.Context) map[string]uint64 {
	resp, err := http.Get(fmt.Sprintf("%s:%d/stats", appConfig.Config.InventoryURI, appConfig.Config.InventoryPort))
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	var stats map[string]uint64

	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	return stats
}

func Scan(ctx context.Context, args []string) (*inventory.Book, error) {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	path := fs.String("barcode-path", "./bad.png", "barcode path")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	code, err := barcode.Read(*path)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(fmt.Sprintf("%s:%d/inventory/%s", appConfig.Config.InventoryURI, appConfig.Config.InventoryPort, code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var book inventory.Book

	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return nil, err
	}

	return &book, nil
}

func Generate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	ids := fs.String("codes", "", "Comma-delimited list of ids")
	path := fs.String("path", ".", "where to store the barcode files")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return barcode.Create(*ids, *path)
}

func CheckOutCheckInBook(ctx context.Context, args []string, cmd string) (*inventory.Book, error) {
	fs := flag.NewFlagSet("checkout", flag.ContinueOnError)
	path := fs.String("barcode-path", "./bad.png", "barcode path")
	who := fs.String("who", "", "who is checking out the book")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if strings.Trim(*who, "") == "" {
		return nil, errors.New("Book cannot be checked out anonymously")
	}
	code, err := barcode.Read(*path)
	if err != nil {
		return nil, err
	}
	req := &types.PostBookRequest{
		Who:  *who,
		Code: code,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("%s:%d/%s", appConfig.Config.InventoryURI, appConfig.Config.InventoryPort, cmd), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var book inventory.Book

	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return nil, err
	}

	return &book, nil
}

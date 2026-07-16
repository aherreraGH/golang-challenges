package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func handleCommands(ctx context.Context, cmd string, cArgs []string) {

	switch cmd {
	case "stats":
		stats := GetInventoryStatistics(ctx)
		if stats != nil {
			fmt.Printf("Statistics: %v", stats)
		}
	case "list":
		books := ListInventory(ctx)

		if books != nil {
			for i, book := range books {
				fmt.Printf("--[%v]-----\n", i)
				fmt.Printf("Title       : %s\n", book.Title)
				fmt.Printf("Author      : %s\n", book.Author)
				fmt.Printf("Barcode ID  : %s\n", book.BarcodeID)
				fmt.Printf("Status      : %s\n", book.Status)
				if strings.ToLower(book.Status) == "not available" {
					fmt.Printf("Checked Out : %s\n", book.CheckedOutAt)
					fmt.Printf("Who         : %s\n", book.CheckedOutBy)
					fmt.Printf("Due         : %s\n", book.DueDate)
				}
				fmt.Println()
			}
		} else {
			fmt.Println("No books found, or a server error occurred.")
		}
	case "scan":
		// .\bin\kiosk.exe scan -barcode-path ".\barcodes\LIB-GO-001-A7.png"
		book, err := Scan(ctx, cArgs)
		if err != nil {
			fmt.Println("Failed to scan:", err.Error())
			os.Exit(1)
		} else {
			if strings.Trim(book.Title, "") == "" {
				fmt.Println("Book not found")
			} else {
				fmt.Printf("Title       : %s\n", book.Title)
				fmt.Printf("Author      : %s\n", book.Author)
				fmt.Printf("Barcode ID  : %s\n", book.BarcodeID)
				fmt.Printf("Status      : %s\n", book.Status)
				if strings.ToLower(book.Status) == "not available" {
					fmt.Printf("Checked Out : %s\n", book.CheckedOutAt)
					fmt.Printf("Who         : %s\n", book.CheckedOutBy)
					fmt.Printf("Due         : %s\n", book.DueDate)
				}
				fmt.Println()
			}
		}
	case "checkout":
		// .\bin\kiosk.exe checkout -who "Joe" -barcode-path ".\barcodes\LIB-GO-001-A7.png"
		book, err := CheckOutCheckInBook(ctx, cArgs, cmd)
		if err != nil {
			fmt.Println("Failed to scan:", err.Error())
			os.Exit(1)
		} else {
			if strings.Trim(book.Title, "") == "" {
				fmt.Println("Book not found")
			} else {
				fmt.Printf("Title       : %s\n", book.Title)
				fmt.Printf("Author      : %s\n", book.Author)
				fmt.Printf("Barcode ID  : %s\n", book.BarcodeID)
				fmt.Printf("Status      : %s\n", book.Status)
				if strings.ToLower(book.Status) == "not available" {
					fmt.Printf("Checked Out : %s\n", book.CheckedOutAt)
					fmt.Printf("Who         : %s\n", book.CheckedOutBy)
					fmt.Printf("Due         : %s\n", book.DueDate)
				}
				fmt.Println()
			}
		}
	case "checkin":
		// .\bin\kiosk.exe checkout -who "Joe" -barcode-path ".\barcodes\LIB-GO-001-A7.png"
		book, err := CheckOutCheckInBook(ctx, cArgs, cmd)
		if err != nil {
			fmt.Println("Failed to scan:", err.Error())
			os.Exit(1)
		} else {
			if strings.Trim(book.Title, "") == "" {
				fmt.Println("Book not found")
			} else {
				fmt.Printf("Title       : %s\n", book.Title)
				fmt.Printf("Author      : %s\n", book.Author)
				fmt.Printf("Barcode ID  : %s\n", book.BarcodeID)
				fmt.Printf("Status      : %s\n", book.Status)
				if strings.ToLower(book.Status) == "not available" {
					fmt.Printf("Checked Out : %s\n", book.CheckedOutAt)
					fmt.Printf("Who         : %s\n", book.CheckedOutBy)
					fmt.Printf("Due         : %s\n", book.DueDate)
				}
				fmt.Printf("Checked In : %s\n", book.CheckedInAt)
				fmt.Printf("Amount Due  : $%.2f\n", book.AmountDue)
				fmt.Println()
			}
		}
	case "generate":
		// .\bin\kiosk.exe generate -codes "LIB-GO-001-A7,LIB-SYS-204-K9,LIB-MATH-88-Z3" -path "./barcodes"
		err := Generate(cArgs)
		if err != nil {
			fmt.Println("Failed to generate:", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("Created barcode image files")
		}
	default:
		fmt.Printf("error: unknown command - %q\n", cmd)
	}

}

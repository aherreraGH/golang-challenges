package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

/**
 * Ensure the command arg[1] is a valid command keyword
 * if not, show some error.
 */
func handleCommands(ctx context.Context, cmd string, cArgs []string) {

	switch cmd {
	case "stats":
		// display a simple statistical output
		stats := GetInventoryStatistics(ctx)
		if stats != nil {
			fmt.Printf("Statistics: %v", stats)
		}
	case "list":
		// show current state of all inventory items
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
		// simulate scanning a barcode - use local image for this challenge.
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
		// simulate checking out a book, use the local barcode images.
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
		// simulate checking in a book, use the local image.
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
		// Generate barcodes based on the code text. This is needed so that the package to read the barcodes
		// can properly convert the generated image to the correct barcode text value.
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

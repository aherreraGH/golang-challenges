package inventory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"localpractice7.com/challenges/internal/checkout"
)

// the book struct object.
type Book struct {
	BarcodeID    string    `json:"barcode_id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Category     string    `json:"category"`
	AcquiredDate time.Time `json:"acquired_date,format:date"`
	Status       string    `json:"status"`
	CheckedOutBy string    `json:"checkedOutBy"`
	CheckedOutAt time.Time `json:"checkedOutAt,format:date"`
	CheckedInAt  time.Time `json:"checkedInAt,format:date"`
	DueDate      time.Time `json:"dueDate,format:date"`
	AmountDue    float64   `json:"amountDue"`
}

// Inventory object, contains 0-many books.
// Also shows the use of sync/mutex, which will use RLock/Lock to ensure deadlocks do not occur.
type Inventory struct {
	mu    sync.RWMutex
	books []*Book
	stats *Statistics
}

func (i *Inventory) Load(ctx context.Context, path string) (string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// check if inventory file exists
	_, err := os.Stat(path)
	if err != nil {
		return "File not found", err
	}
	// load inventory file
	data, err := os.ReadFile(path)
	if err != nil {
		return "Failed to read file", err
	}

	err = json.Unmarshal(data, &i.books)
	if err != nil {
		return "Failed to convert to books", err
	}

	i.stats = new(Statistics)

	return "Inventory loaded", nil
}

func (i *Inventory) Save(ctx context.Context) error {
	return nil
}

func (i *Inventory) Find(code string) (*Book, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	i.stats.IncRequests()
	i.stats.IncSearch()
	for _, item := range i.books {
		if item.BarcodeID == code {
			return cloneBook(item), nil
		}
	}
	book := new(Book)

	i.stats.IncSearch()

	return book, errors.New("Not found")
}

func (i *Inventory) List(ctx context.Context) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	i.stats.IncRequests()
	for _, item := range i.books {
		fmt.Println("Item: ", item)
	}
}

func (i *Inventory) Get() []*Book {
	i.mu.RLock()
	defer i.mu.RUnlock()

	i.stats.IncRequests()

	result := make([]*Book, len(i.books))

	for idx, book := range i.books {
		copy := *book
		result[idx] = &copy
	}

	return result
}

func (i *Inventory) GetStatistics() map[string]uint64 {
	i.stats.IncRequests()
	return i.stats.DisplayStats()
}

func (i *Inventory) CheckOut(code string, who string) (*Book, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.stats.IncRequests()
	i.stats.IncCheckouts()
	// find book
	book, err := i.findUnsafe(code)
	if err != nil {
		return book, err
	}
	// is book already checked out?
	if strings.ToLower(book.Status) != "available" {
		return book, errors.New("Book is not avaiable")
	}

	book.CheckedOutAt = time.Now()
	book.CheckedOutBy = who
	book.DueDate = checkout.WhenDue(book.CheckedOutAt)
	book.Status = "Not available"

	return book, nil
}

func (i *Inventory) CheckIn(code string) (*Book, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.stats.IncRequests()
	i.stats.IncCheckins()
	// find book
	book, err := i.findUnsafe(code)
	if err != nil {
		return book, err
	}
	book.CheckedInAt = time.Now()
	book.AmountDue = checkout.AmountDue(book.CheckedInAt)
	book.Status = "Available"
	book.CheckedOutBy = ""
	book.CheckedOutAt = time.Time{}

	return book, nil
}

func (i *Inventory) findUnsafe(code string) (*Book, error) {
	for _, item := range i.books {
		if item.BarcodeID == code {
			return item, nil
		}
	}

	return nil, errors.New("Not found")
}

func cloneBook(b *Book) *Book {
	copy := *b
	return &copy
}

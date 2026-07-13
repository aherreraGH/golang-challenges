package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"localpractice6.com/challenges/types"
)

var (
	fileName string = "./task.json"
)

type Store struct {
	mu    sync.RWMutex
	tasks []types.Task

	/**
	A good rule of thumb

	Use a mutex if protecting a complex object?

	map
	slice
	struct
	linked list
	cache

	Use atomic if just one value?

	bool
	int
	uint64
	pointer
	*/

	reads  atomic.Uint64
	writes atomic.Uint64
	saves  atomic.Uint64
	dirty  atomic.Bool
}

func (s *Store) Delete(ctx context.Context) error {
	err := os.Remove(fileName) // replace "test.txt" with your file name
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Load(ctx context.Context) error {

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		// create it
		err = s.Save(ctx)
		if err != nil {
			return err
		}
	}

	// Open the JSON file
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer jsonFile.Close()

	// Decode the JSON data
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&s.tasks); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	return nil
}
func (s *Store) Add(ctx context.Context, title string) (types.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task := types.Task{
		ID:        len(s.tasks),
		Title:     title,
		Completed: false,
	}
	s.tasks = append(s.tasks, task)
	// fmt.Println("Added task...", s.tasks)
	return task, nil
}
func (s *Store) Complete(ctx context.Context, id int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	found := false
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Completed = true
			found = true
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("Task ID %v not found", id))
	}
	return nil
}
func (s *Store) List(ctx context.Context) ([]types.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]types.Task, len(s.tasks))
	copy(result, s.tasks)

	return result, nil
}

func (s *Store) Save(ctx context.Context) error {
	// ensure file exists
	f, _ := os.Create(fileName)
	defer f.Close()

	dataToSave, _ := json.MarshalIndent(s.tasks, "", "\t")

	_, err := f.Write(dataToSave)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) RunAutoSave(ctx context.Context, interval time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Auto save triggered...")
			if err := s.Save(ctx); err != nil &&
				!errors.Is(err, context.Canceled) {
				fmt.Printf("autosave failed: %v\n", err)
			}
		case <-ctx.Done():
			fmt.Println("autosave stopped")
			return
		}
	}
}

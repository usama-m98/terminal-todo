package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Task struct {
	Title  string
	Status string
}

type TaskStore struct {
	store map[int]Task
}

func NewTaskStore() (*TaskStore, error) {
	ts := &TaskStore{}
	store, err := ts.loadDB()
	if err != nil {
		return &TaskStore{}, err
	}
	ts.store = store
	return ts, nil
}

func (s *TaskStore) addToStore(task string) (int, error) {
	key := 0
	for k := range s.store {
		// overwrites till it's found
		key = k + 1
	}
	s.store[key] = Task{
		Title:  task,
		Status: "not complete",
	}

	err := s.writeToDB()

	if err != nil {
		return 0, fmt.Errorf("failed to add to store: %w", err)
	}

	return key, nil
}

func (s *TaskStore) markComplete(id int) error {
	if entry, ok := s.store[id]; !ok {
		return errors.New("ID does not exist")
	} else {
		entry.Status = "completed"
		s.store[id] = entry
	}

	return nil
}

func (s *TaskStore) writeToDB() error {
	dat, err := json.Marshal(s.store)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	err = os.WriteFile("./task-store.json", dat, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to task-store.json %w", err)
	}

	s.loadDB()
	return nil
}

func (s *TaskStore) removeFromStore(id int) error {
	delete(s.store, id)

	err := s.writeToDB()
	if err != nil {
		return fmt.Errorf("failed to remove from store: %w", err)
	}
	return nil
}

func (s *TaskStore) loadDB() (map[int]Task, error) {
	err := s.createDB()
	if err != nil {
		return map[int]Task{}, err
	}
	dat, err := os.ReadFile("./task-store.json")
	if err != nil {
		return map[int]Task{}, fmt.Errorf("error reading task-store.json %w", err)
	}

	var tasksDat map[int]Task
	err = json.Unmarshal(dat, &tasksDat)
	if err != nil {
		return map[int]Task{}, fmt.Errorf("failed to load db %w", err)
	}

	return tasksDat, nil
}

// create json file if it doesn't exist
func (s *TaskStore) createDB() error {
	if _, err := os.Stat("./task-store.json"); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile("task-store.json", []byte("{}"), 0644)
		if err != nil {
			return fmt.Errorf("Error creating db: %w", err)
		}
	}

	return nil
}

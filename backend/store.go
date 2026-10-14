package main

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// Todo is a single to-do item.
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

var ErrNotFound = errors.New("todo not found")

// Store is an in-memory, concurrency-safe todo store.
// Data is lost when the process exits (by design, for learning purposes).
type Store struct {
	mu     sync.RWMutex
	nextID int
	todos  map[int]Todo
}

func NewStore() *Store {
	return &Store{nextID: 1, todos: make(map[int]Todo)}
}

func (s *Store) List() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Todo, 0, len(s.todos))
	for _, t := range s.todos {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *Store) Get(id int) (Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.todos[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	return t, nil
}

func (s *Store) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := Todo{ID: s.nextID, Title: title, CreatedAt: time.Now().UTC()}
	s.todos[t.ID] = t
	s.nextID++
	return t
}

func (s *Store) Update(id int, title string, done bool) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	t.Title = title
	t.Done = done
	s.todos[id] = t
	return t, nil
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return nil
}

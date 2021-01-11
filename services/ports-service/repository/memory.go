package repository

import (
	"context"
	"fmt"
)

//I would have also used a cache implementation (like "github.com/patrickmn/go-cache"), but for simplicity I'll go with map

type MemoryStorage map[string]string

func NewMemoryStorage() MemoryStorage {
	return make(MemoryStorage)
}

func (r MemoryStorage) Get(ctx context.Context, id string) (string, *RepoError) {
	val, ok := r[id]
	if !ok {
		return "", &RepoError{NOT_FOUND, fmt.Errorf("ID %s not found", id)}
	}
	return val, nil
}

func (r MemoryStorage) Set(ctx context.Context, id, value string) *RepoError {
	r[id] = value
	return nil
}

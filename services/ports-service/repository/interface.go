package repository

import "context"

type ErrorType int

const (
	// could have used iota, but prefer it with assigned values for visibility and contro
	NONE           ErrorType = 0
	NOT_FOUND      ErrorType = 1
	ALREADY_EXISTS ErrorType = 2
	INTERNAL       ErrorType = 3
	UNKNOWN        ErrorType = 4
)

type RepoError struct {
	Type ErrorType
	Err  error
}

type Repository interface {
	Get(ctx context.Context, id string) (string, *RepoError)
	Set(ctx context.Context, id, value string) *RepoError
}

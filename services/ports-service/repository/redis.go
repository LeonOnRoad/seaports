package repository

import (
	"context"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(c *redis.Client) *Redis {
	return &Redis{
		client: c,
	}
}

func (r Redis) Get(ctx context.Context, id string) (string, *RepoError) {

	val, err := r.client.Get(id).Result()
	if err != nil {
		if err == redis.Nil {
			return "", &RepoError{Type: NOT_FOUND, Err: err}
		}
		return "", &RepoError{Type: INTERNAL, Err: err}
	}

	return val, nil
}

func (r Redis) Set(ctx context.Context, id, value string) *RepoError {
	err := r.client.Set(id, value, 0).Err()
	if err != nil {
		return &RepoError{Type: INTERNAL, Err: err}
	}
	return nil
}

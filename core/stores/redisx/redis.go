package redisx

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addrs    string
	Password string
	Index    int
}

func NewNodeRds(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addrs,
		Password: cfg.Password,
		DB:       cfg.Index,
	})
}

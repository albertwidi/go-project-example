package redis

import (
	"context"
	"errors"
)

// error list
var (
	ErrResponseNotOK = errors.New("redis: response is not ok")
)

// Redis interface
type Redis interface {
	Set(ctx context.Context, key string, value interface{}) (string, error)
	SetNX(ctx context.Context, key string, value interface{}, expire int) (int, error)
	SetEX(ctx context.Context, key string, value interface{}, expire int) (string, error)
	Get(ctx context.Context, key string) (string, error)
	MSet(ctx context.Context, pairs ...interface{}) (string, error)
	MGet(ctx context.Context, keys ...string) ([]string, error)
	HSet(ctx context.Context, key, field string, value interface{}) (int, error)
	HSetEX(ctx context.Context, key, field string, value interface{}, expire int) (int, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HMSet(ctx context.Context, key string, kv map[string]interface{}) (string, error)
	HMGet(ctx context.Context, key string, fields ...string) ([]string, error)
	HDel(ctx context.Context, key string, fields ...string) (int, error)
	LLen(ctx context.Context, key string) (int, error)
	LIndex(ctx context.Context, key string, index int) (string, error)
	LSet(ctx context.Context, key, value string, index int) (int, error)
	LPush(ctx context.Context, key string, values ...interface{}) (int, error)
	LPushX(ctx context.Context, key string, values ...interface{}) (int, error)
	LPop(ctx context.Context, key string) (string, error)
	LRem(ctx context.Context, key, value string, count int) (int, error)
	LTrim(ctd context.Context, key string, start, stop int) (string, error)
}

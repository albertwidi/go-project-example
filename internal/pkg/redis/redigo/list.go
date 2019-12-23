package redigo

import (
	"context"

	"github.com/albertwidi/go-project-example/internal/pkg/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// LLen get the length of the list
func (rdg *Redigo) LLen(ctx context.Context, key string) (int, error) {
	result, err := redigo.Int(rdg.do(ctx, redis.CommandLLen, key))
	if err != nil && !rdg.IsErrNil(err) {
		return 0, err
	}
	return result, err
}

// LIndex to get value from a certain list index
func (rdg *Redigo) LIndex(ctx context.Context, key string, index int) (string, error) {
	result, err := redigo.String(rdg.do(ctx, redis.CommandLIndex, key, index))
	if err != nil && !rdg.IsErrNil(err) {
		return "", err
	}
	return result, err
}

// LSet to set value to some index
func (rdg *Redigo) LSet(ctx context.Context, key, value string, index int) (int, error) {
	result, err := redigo.Int(rdg.do(ctx, redis.CommandLSET, index, value))
	if err != nil && !rdg.IsErrNil(err) {
		return 0, err
	}
	return result, err
}

// LPush prepend values to the list
func (rdg *Redigo) LPush(ctx context.Context, key string, values ...interface{}) (int, error) {
	args := make([]interface{}, len(values)+1)
	args[0] = key
	for i, value := range values {
		args[i+1] = value
	}

	result, err := redigo.Int(rdg.do(ctx, redis.CommandLPush, args...))
	if err != nil && !rdg.IsErrNil(err) {
		return 0, err
	}
	return result, err
}

// LPushX prepend values to the list
func (rdg *Redigo) LPushX(ctx context.Context, key string, values ...interface{}) (int, error) {
	args := make([]interface{}, len(values)+1)
	args[0] = key
	for i, value := range values {
		args[i+1] = value
	}

	result, err := redigo.Int(rdg.do(ctx, redis.CommandLPushX, args...))
	if err != nil && !rdg.IsErrNil(err) {
		return 0, err
	}
	return result, err
}

// LPop removes and get the first element in the list
func (rdg *Redigo) LPop(ctx context.Context, key string) (string, error) {
	result, err := redigo.String(rdg.do(ctx, redis.CommandLPop, key))
	if err != nil && !rdg.IsErrNil(err) {
		return "", err
	}
	return result, err
}

// LRem command
func (rdg *Redigo) LRem(ctx context.Context, key, value string, count int) (int, error) {
	result, err := redigo.Int(rdg.do(ctx, redis.CommandLRem, key, count, value))
	if err != nil && !rdg.IsErrNil(err) {
		return 0, err
	}
	return result, err
}

// LTrim command
func (rdg *Redigo) LTrim(ctx context.Context, key string, start, stop int) (string, error) {
	result, err := redigo.String(rdg.do(ctx, redis.CommandLTrim, key, start, stop))
	if err != nil && !rdg.IsErrNil(err) {
		return "", err
	}
	return result, err
}

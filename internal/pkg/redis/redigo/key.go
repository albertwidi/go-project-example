package redigo

import (
	"context"

	"github.com/albertwidi/go_project_example/internal/pkg/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// Set key and value
func (rdg *Redigo) Set(ctx context.Context, key string, value interface{}) (string, error) {
	ok, err := redigo.String(rdg.do(ctx, "SET", key, value))
	if !IsResponseOK(ok) {
		return ok, redis.ErrResponseNotOK
	}
	return ok, err
}

// SetNX do SETNX (only set if not exist) with SET's NX & EX args.
// It sets the key which will expired in `expire` seconds
func (rdg *Redigo) SetNX(ctx context.Context, key string, value interface{}, expire int) (int, error) {
	resp, err := redigo.Int(rdg.do(ctx, "SETNX", key, value, "NX", "EX", expire))
	if err != nil && !IsErrNil(err) {
		return 0, err
	}
	return resp, err
}

// SetEX key and value
// It sets the key wich will expired in `expire` seconds
func (rdg *Redigo) SetEX(ctx context.Context, key string, value interface{}, expire int) (string, error) {
	resp, err := redigo.String(rdg.do(ctx, "SETEX", key, expire, value))
	if err != nil && !IsErrNil(err) {
		return "", err
	}
	return resp, err
}

// Get string value
func (rdg *Redigo) Get(ctx context.Context, key string) (string, error) {
	resp, err := redigo.String(rdg.do(ctx, "GET", key))
	if err != nil && !IsErrNil(err) {
		return "", err
	}
	return resp, err
}

// MSet keys and values
// please use basic types only (no struct, array, or map) for arguments
func (rdg *Redigo) MSet(ctx context.Context, pairs ...interface{}) (string, error) {
	ok, err := redigo.String(rdg.do(ctx, "MSET", pairs...))
	if !IsResponseOK(ok) {
		return ok, redis.ErrResponseNotOK
	}
	return ok, err
}

// MGet keys
func (rdg *Redigo) MGet(ctx context.Context, keys ...string) ([]string, error) {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}

	resp, err := redigo.Strings(rdg.do(ctx, "MGET", args...))
	if err != nil && !IsErrNil(err) {
		return nil, err
	}
	return resp, err
}

// HSet field and value based on key
func (rdg *Redigo) HSet(ctx context.Context, key, field string, value interface{}) (int, error) {
	resp, err := redigo.Int(rdg.do(ctx, "HSET", key, field, value))
	if err != nil && !IsErrNil(err) {
		return resp, err
	}
	return resp, err
}

// HSetEX key and value and sets the expiration to the given `expire` seconds
func (rdg *Redigo) HSetEX(ctx context.Context, key, field string, value interface{}, expire int) (int, error) {
	conn, err := rdg.getConn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	resp, err := redigo.Int(conn.Do("HSET", key, field, value))
	if err != nil && !IsErrNil(err) {
		return resp, err
	}

	resp, err = redigo.Int(rdg.do(ctx, "EXPIRE", key, expire))
	if err != nil && !IsErrNil(err) {
		return resp, err
	}

	return resp, err
}

// HGet key and value
func (rdg *Redigo) HGet(ctx context.Context, key, field string) (string, error) {
	resp, err := redigo.String(rdg.do(ctx, "HGET", key, field))
	if err != nil && !IsErrNil(err) {
		return resp, err
	}
	return resp, err
}

// HGetAll key and value
func (rdg *Redigo) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	resp, err := redigo.Strings(rdg.do(ctx, "HGETALL", key))
	if err != nil && !IsErrNil(err) {
		return nil, err
	}

	kv := make(map[string]string)
	respLen := len(resp)

	for i := 0; i < respLen; i += 2 {
		kv[resp[i]] = resp[i+1]
	}
	return kv, err
}

// HMSet function
// please use basic types only (no struct, array, or map) for kv value
func (rdg *Redigo) HMSet(ctx context.Context, key string, kv map[string]interface{}) (string, error) {
	var (
		args = make([]interface{}, 1+(len(kv)*2))
		idx  = 1
	)
	args[0] = key
	for k, v := range kv {
		args[idx] = k
		args[idx+1] = v
		idx += 2
	}

	resp, err := redigo.String(rdg.do(ctx, "HMSET", args...))
	if err != nil && !IsErrNil(err) {
		return resp, err
	}
	return resp, err
}

// HMGet keys and value
func (rdg *Redigo) HMGet(ctx context.Context, key string, fields ...string) ([]string, error) {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}

	resp, err := redigo.Strings(rdg.do(ctx, "HMGET", args...))
	if err != nil && !IsErrNil(err) {
		return resp, err
	}
	return resp, err
}

// HDel fields of a key
func (rdg *Redigo) HDel(ctx context.Context, key string, fields ...string) (int, error) {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}

	resp, err := redigo.Int(rdg.do(ctx, "HDEL", args...))
	if err != nil && !IsErrNil(err) {
		return resp, err
	}
	return resp, err
}

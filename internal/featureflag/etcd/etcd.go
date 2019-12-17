package etcd

import (
	"context"
	"errors"
	"time"

	"go.etcd.io/etcd/client"
)

// Etcd backend for feature flag
type Etcd struct {
	c       client.Client
	keysAPI client.KeysAPI
}

// New etcd feature-flag backend
func New(endpoints []string, timeout time.Duration) (*Etcd, error) {
	c, err := client.New(client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: timeout,
	})
	if err != nil {
		return nil, err
	}

	etcd := Etcd{
		c:       c,
		keysAPI: client.NewKeysAPI(c),
	}
	return &etcd, nil
}

// Set key to etcd
func (etcd *Etcd) Set(ctx context.Context, key, value string) error {
	resp, err := etcd.keysAPI.Set(ctx, key, value, nil)
	if err != nil {
		return err
	}

	if resp.Index == 0 {
		return errors.New("etcd: set failed, have no index")
	}
	return nil
}

// Get key from etcd
func (etcd *Etcd) Get(ctx context.Context, key string) (string, error) {
	resp, err := etcd.keysAPI.Get(ctx, key, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

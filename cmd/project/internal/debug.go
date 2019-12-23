package project

import (
	debugserver "github.com/albertwidi/go-project-example/internal/server/debug"
)

func newDebugServer(address string, r *Repositories) (*debugserver.Server, error) {
	usecases := debugserver.Usecases{}
	s, err := debugserver.New(address, usecases)
	if err != nil {
		return nil, err
	}

	return s, nil
}

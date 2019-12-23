package main

import (
	"errors"
	"log"

	"github.com/albertwidi/go-project-example/internal/xerrors"
)

var (
	errSomething = errors.New("something is error bruh")
)

func main() {
	// xerrors.SetCaller(true)
	err := a()
	log.Printf("%v", err)

	if xerrors.Is(err, errSomething) {
		log.Println("WAIKI")
	} else {
		log.Println("WAHH")
	}

	e := xerrors.XUnwrap(err)
	log.Println(e.Kind())
}

func a() error {
	err := b()
	return xerrors.New(xerrors.Op("function/a"), err)
}

func b() error {
	err := c()
	return xerrors.New(xerrors.Op("function/b"), err)
}

func c() error {
	err := d()
	return xerrors.New(xerrors.Op("function/c"), err)
}

func d() error {
	return xerrors.New(xerrors.Op("function/d"), errSomething, xerrors.KindBadRequest)
}

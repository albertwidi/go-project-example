package server

import (
	"io/ioutil"
	"net"

	httpapi "github.com/albertwidi/kothak/api/http"
	orderresource "github.com/albertwidi/kothak/resource/order"
	paymentresource "github.com/albertwidi/kothak/resource/payment"
	userresource "github.com/albertwidi/kothak/resource/user"
	orderservice "github.com/albertwidi/kothak/service/order"
	paymentservice "github.com/albertwidi/kothak/service/payment"
	userservice "github.com/albertwidi/kothak/service/user"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

// Main program or run the server
func Main() error {
	// read config from config directory
	out, err := ioutil.ReadFile("config/kothak.config.yml")
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(out, &Config); err != nil {
		return err
	}

	masterDB, err := sqlx.Connect("postgres", Config.Database.Master)
	if err != nil {
		return err
	}
	followerDB, err := sqlx.Connect("postgres", Config.Database.Follower)
	if err != nil {
		return err
	}

	// user
	userres := userresource.New(masterDB, followerDB)
	usersvc := userservice.New(userres)
	// order
	orderres := orderresource.New(masterDB, followerDB)
	ordersvc := orderservice.New(orderres, usersvc)
	// payment
	paymentres := paymentresource.New(masterDB, followerDB)
	paymentsvc := paymentservice.New(paymentres, usersvc)

	// create a new listener for http and grpc server
	listener, err := net.Listen("tcp", "8000")
	if err != nil {
		return err
	}

	httpserver := httpapi.Server{
		UserService:    usersvc,
		OrderService:   ordersvc,
		PaymentService: paymentsvc,
	}
	return httpserver.Serve(listener)
}

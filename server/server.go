package server

import (
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	orderresource "gitlab.com/kosanapp/kothak/resource/order"
	userresource "gitlab.com/kosanapp/kothak/resource/user"
	orderservice "gitlab.com/kosanapp/kothak/services/order"
	userservice "gitlab.com/kosanapp/kothak/services/user"
	"gopkg.in/yaml.v2"
)

// Execute or run the server
func Main() error {
	// read config from config directory
	out, err := ioutil.ReadFile("../config/kothak.config.yml")
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

	userres := userresource.New(masterDB, followerDB)
	usersvc := userservice.New(userres)

	orderres := orderresource.New(masterDB, followerDB)
	ordersvc := orderservice.New(orderres, usersvc)
	return nil
}

package user

import (
	"net/http"

	"github.com/albertwidi/go_project_example/api"
)

var userService api.UserService

func Init(service api.UserService) {
	userService = service
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from user service"))
}

func IsUserActive(w http.ResponseWriter, r *http.Request) {
}

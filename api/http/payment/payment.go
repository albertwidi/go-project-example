package payment

import (
	"net/http"

	"github.com/albertwidi/kothak/api"
)

var paymentService api.PaymentService

func Init(service api.PaymentService) {
	paymentService = service
}

func CreatePayment(w http.ResponseWriter, r *http.Request) {

}

func ConfirmPayment(w http.ResponseWriter, r *http.Request) {

}

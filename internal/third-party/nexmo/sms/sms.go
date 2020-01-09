package sms

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/albertwidi/go-project-example/internal/pkg/http/request"
)

// Client sms module for nexmo
type Client struct {
	httpClient *http.Client
	config     Config
}

// Config of nexmo
type Config struct {
	APIKey           string
	APISecret        string
	Endpoint         string
	CallbackEndpoint string
}

// Validate nexmo config
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("api-key is needed")
	}

	if c.APISecret == "" {
		return errors.New("api-secret is needed")
	}

	if c.Endpoint == "" {
		// use default endpoint
		c.Endpoint = "https://rest.nexmo.com/sms/json"
	}

	return nil
}

// New nexmo sms module
func New(config Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	client := Client{
		httpClient: new(http.Client),
		config:     config,
	}

	return &client, nil
}

// Payload of sms
type Payload struct {
	From    string
	To      string
	Message string
}

// Request data
type Request struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	// hash of request parameter
	Sig  string `json:"sig"`
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
	// NOT REQUIRED, nexmo attempt to delivery
	TTL int `json:"ttl"`
	// status report by nexmo
	StatusReportReq bool `json:"status-report-req"`
	// where the callback is going to, nexmo will send back the status of send
	Callback string `json:"callback"`
	Type     string `json:"type"`
}

// Response struct
type Response struct {
	MessageCount string `json:"message-count"`
	Messages     []struct {
		To               string `json:"to"`
		MessageID        string `json:"message_id"`
		Status           string `json:"status"`
		RemainingBalance string `json:"remaining_balance"`
		MessagePrice     string `json:"message-price"`
		Network          string `json:"network"`
		ErrorText        string `json:"error-text"`
	} `json:"messages"`
}

// NexmoSMSCallback data
type NexmoSMSCallback struct {
	// the number that message sent to
	MSISDN string `json:"msisdn"`
	// this is the sender_id
	To          string `json:"to"`
	NetworkCode string `json:"network_code"`
	// message_id from nexmo
	MessageID string `json:"message_id"`
	// price of the message
	Price string `json:"price"`
	// status of delivery
	Status string `json:"status"`
	Scts   string `json:"scts"`
	// should be not 0 if error
	ErrCode string `json:"err-code"`
	// date of webhook triggered
	MessageTimestamp string `json:"message_timestamp"`
}

// Send sms using nexmo
// currently, the API only expect to send 1 message
func (c *Client) Send(ctx context.Context, payload Payload) (Response, error) {
	httpreq, err := request.New(ctx).
		Post(c.config.Endpoint).
		PostForm("api_key", c.config.APIKey,
			"api_secret", c.config.APISecret,
			"from", payload.From,
			"to", payload.To,
			"text", payload.Message).
		Headers(request.Header().ContentType().ApplicationFormWWWURLEncoded().Headers()).
		Compile()
	if err != nil {
		return Response{}, err
	}

	resp, err := c.httpClient.Do(httpreq)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	apiResp := Response{}
	if err := json.Unmarshal(out, &apiResp); err != nil {
		return Response{}, err
	}

	if resp.StatusCode >= 300 {
		err := errors.New("failed to send otp")
		return Response{}, err
	}

	messageCount, err := strconv.Atoi(apiResp.MessageCount)
	if err != nil {
		return Response{}, err
	}

	// no message being sent
	if messageCount == 0 {
		return Response{}, errors.New("nexmo: no message sent")
	}

	message := apiResp.Messages[0]
	if message.Status != "0" {
		return Response{}, errors.New(message.ErrorText)
	}

	return apiResp, nil
}

// Callback to handle callback from nexmo
func (c *Client) Callback() {

}

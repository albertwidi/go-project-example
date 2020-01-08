package pushmessage

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

// Firebase backend for pushmessage
type Firebase struct {
	msgclient *messaging.Client
}

// Config of firebase
type Config struct {
	ProjectID          string
	ServiceAccountID   string
	Bucket             string
	ServiceAccountFile string
	DryRun             bool
}

// SendOptions for firebase client
type SendOptions struct {
	DryRun bool
}

// New firebase push message package
func New(ctx context.Context, config *Config) (*Firebase, error) {
	var opts option.ClientOption
	cfg := new(firebase.Config)
	if config != nil {
		cfg.ProjectID = config.ProjectID
		cfg.ServiceAccountID = config.ServiceAccountID

		if config.ServiceAccountFile != "" {
			opts = option.WithCredentialsFile(config.ServiceAccountFile)
		}
	}

	app, err := firebase.NewApp(ctx, cfg, opts)
	if err != nil {
		return nil, err
	}

	msgclient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	f := Firebase{
		msgclient: msgclient,
	}

	return &f, nil
}

// Send notification
func (f *Firebase) Send(ctx context.Context, message *messaging.Message, options *SendOptions) (string, error) {
	var opts SendOptions
	if options != nil {
		opts = *options
	}

	var (
		id  string
		err error
	)

	if !opts.DryRun {
		id, err = f.msgclient.Send(ctx, message)
	} else {
		// use default token if token not exists
		if message.Token == "" {
			message.Token = dummyToken
		}
		id, err = f.msgclient.SendDryRun(ctx, message)
	}
	return id, err
}

// the oauth package implements google oauth package

package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// scope list
const (
	ScopeEmail   = "email"
	ScopeProfile = "profile"
)

// OAuth struct
type OAuth struct {
	config *oauth2.Config
}

// Config struct
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// New oauth
func New(config Config) *OAuth {
	oauth := OAuth{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes:       config.Scopes,
			Endpoint:     google.Endpoint,
		},
	}
	return &oauth
}

// Config return google oauth2 configuration
func (oauth *OAuth) Config() *oauth2.Config {
	return oauth.config
}

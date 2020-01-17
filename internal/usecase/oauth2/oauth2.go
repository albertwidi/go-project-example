package oauth2

import "errors"

var (
	allowedGrantType = map[string]bool{
		"password":           true,
		"client_credentials": true,
	}
)

// Usecase for oauth2
type Usecase struct {
}

// New oauth2 usecase
func New() *Usecase {
	u := Usecase{}
	return &u
}

// IsGrantTypeAllowed to check whether the grant type is allowed or not
func (u *Usecase) IsGrantTypeAllowed(grantType string) error {
	_, ok := allowedGrantType[grantType]
	if !ok {
		return errors.New("oauth2: grant type is not allowed")
	}
	return nil
}

// Create new oauth2 access
func (u *Usecase) Create(userID string, scopes []string) error {
	return nil
}

// Authorize user via oauth2
func (u *Usecase) Authorize(grantType string, authParam interface{}) error {
	return nil
}

// AuthGrantPassword to authorize user via
type AuthGrantPassword struct {
	Username string
	Password string
}

// authorizeGrantPassword for authorizing user via password
// authorize with password grant always return the primary oauth2 token that user has
// this grant type is used in user login via password
// password can be user encrypted password in database
// or an otp via sms
func (u *Usecase) authorizeGrantPassword(username, password string) error {
	return nil
}

// authorizeGrantClientCredentials for athorizing user via client credentials(client_id/client_secret)
func (u *Usecase) authorizeGrantClientCredentials(clientID, clientSecret string) error {
	return nil
}
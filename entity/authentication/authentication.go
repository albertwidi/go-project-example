package authentication

// Action of authentication
type Action int

// Provider of authentication
type Provider string

// Authentication entity
type Authentication struct {
	Action   Action
	Provider Provider
	Username string
}

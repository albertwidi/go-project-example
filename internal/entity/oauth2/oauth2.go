package oauth2

// Scope type for oauth2
type Scope int

// list of scope
const (
	ScopeRead Scope = iota + 1
	ScopeWrite
)

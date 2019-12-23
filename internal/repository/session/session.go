// session is stored in a single hash key inside of redis k/v
// for example we have hash of user:{hashed_value}
// key 'data' is reserved for user data information
// and the other key is {session_id}

package session

import "context"

// Repository struct
type Repository struct {
}

// New session repository
func New() {

}

// Create a new session
func (r *Repository) Create(ctx context.Context) error {
	return nil
}

// SaveUserInfo is a special event
// whether user is a new user login, or user is updating the user information
func (r *Repository) SaveUserInfo(ctx context.Context) error {
	return nil
}

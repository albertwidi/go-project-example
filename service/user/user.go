package user

import "context"

type User struct {
	ID int64
}

func (s *Service) IsUserActive(ctx context.Context, userid int64) (bool, error) {
	return true, nil
}

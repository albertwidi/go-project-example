package user

type User struct {
	ID int64
}

func (s *Service) IsUserActive(userid int64) (bool, error) {
	return true, nil
}

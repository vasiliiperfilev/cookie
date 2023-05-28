package data

type StubUserModel struct {
	users   []User
	idCount int64
}

func (s *StubUserModel) Insert(user *User) error {
	if _, err := s.GetByEmail(user.Email); err == nil {
		return ErrDuplicateEmail
	}
	s.idCount++
	user.Id = s.idCount
	s.users = append(s.users, *user)
	return nil
}

func (s *StubUserModel) GetByEmail(email string) (*User, error) {
	for _, user := range s.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, ErrRecordNotFound
}

func (s *StubUserModel) Update(user *User) error {
	if _, err := s.GetByEmail(user.Email); err == nil {
		return ErrDuplicateEmail
	}
	for k, v := range s.users {
		if v.Id == user.Id {
			s.users[k] = *user
		}
	}
	return nil
}

func NewStubUserModel(users []User) *StubUserModel {
	return &StubUserModel{users: users}
}
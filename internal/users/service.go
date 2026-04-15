package users

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, input *User) (*User, error) {
	user, err := NewUser(
		input.Name,
		input.LastName,
		input.Email,
		Role(input.Role),
	)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) TerminateEmployee(ctx context.Context, id int64) (*User, error) {
	return s.repo.TerminateEmployee(ctx, id)
}

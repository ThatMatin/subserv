package service

import (
	"context"
	"fmt"

	"github.com/thatmatin/subserv/internal/repo"
)

type UserService interface {
	Exists(context.Context, uint) (bool, error)
}

type userService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Exists(ctx context.Context, ID uint) (bool, error) {
	exists, err := s.repo.Exists(ctx, ID)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

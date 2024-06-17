package services

import (
	"context"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/repository"
)

func (s *Service) RegisterNewAccount(ctx context.Context, form models.RegisterForm) (*repository.Token, error) {
	return s.accountRepo.RegisterNewAccount(ctx, form)
}

func (s *Service) Login(ctx context.Context, form models.LoginForm) (*repository.Token, error) {
	return s.accountRepo.Login(ctx, form)
}

func (s *Service) Logout(ctx context.Context, token repository.Token) error {
	return s.accountRepo.Logout(ctx, token)
}

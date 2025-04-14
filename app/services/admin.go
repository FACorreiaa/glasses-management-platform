package services

import (
	"context"
	"math"

	"github.com/google/uuid"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/repository"
)

func (s *Service) GetUsers(ctx context.Context, page, pageSize int, orderBy, sortBy, email string) ([]models.UserSession, error) {
	return s.adminRepo.GetUsers(ctx, page, pageSize, orderBy, sortBy, email)
}

func (s *Service) GetUsersByID(ctx context.Context, userID uuid.UUID) (*models.UserSession, error) {
	return s.adminRepo.GetUsersByID(ctx, userID)
}

func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.adminRepo.DeleteUser(ctx, userID)
}

func (s *Service) UpdateUser(ctx context.Context, user models.UpdateUserForm) error {
	return s.adminRepo.UpdateUser(ctx, user)
}

func (s *Service) InsertUser(ctx context.Context, form models.RegisterFormValues) (*repository.Token, error) {
	return s.adminRepo.InsertUser(ctx, form)
}

func (s *Service) GetUsersSum() (int, error) {
	total, err := s.adminRepo.GetUsersSum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (s *Service) GetAdminID(ctx context.Context, userID uuid.UUID) (*models.UserSession, error) {
	return s.adminRepo.GetAdminID(ctx, userID)
}

func (s *Service) GetAdminGlassesDetails(ctx context.Context, page, pageSize int,
	orderBy, sortBy, username string, leftEye, rightEye *float64) ([]models.Glasses, error) {
	return s.glassesRepo.GetAdminGlassesDetails(ctx, page, pageSize, orderBy, sortBy, username, leftEye, rightEye)
}

func (s *Service) GetEmail(ctx context.Context, email string) error {
	return s.adminRepo.GetEmail(ctx, email)
}

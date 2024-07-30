package services

import (
	"context"
	"math"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/google/uuid"
)

func (s *Service) GetGlasses(ctx context.Context, page, pageSize int,
	orderBy, sortBy, reference string, leftEye, rightEye *float64) ([]models.Glasses, error) {
	return s.glassesRepo.GetGlasses(ctx, page, pageSize, orderBy, sortBy, reference, leftEye, rightEye)
}

func (s *Service) GetGlassesByID(ctx context.Context, glassesID uuid.UUID) (*models.Glasses, error) {
	return s.glassesRepo.GetGlassesByID(ctx, glassesID)
}

func (s *Service) DeleteGlasses(ctx context.Context, glassesID uuid.UUID) error {
	return s.glassesRepo.DeleteGlasses(ctx, glassesID)
}

func (s *Service) UpdateGlasses(ctx context.Context, g models.GlassesForm) error {
	return s.glassesRepo.UpdateGlasses(ctx, g)
}

func (s *Service) InsertGlasses(ctx context.Context, g models.GlassesForm) error {
	return s.glassesRepo.InsertGlasses(ctx, g)
}

func (s *Service) GetGlassesByType(ctx context.Context, page, pageSize int,
	orderBy, sortBy, glassesType string) ([]models.Glasses, error) {
	return s.glassesRepo.GetGlassesByType(ctx, page, pageSize, orderBy, sortBy, glassesType)
}

// CALCULATE GET SUM
func (s *Service) GetGlassesByStock(ctx context.Context,
	page, pageSize int, orderBy, sortBy string, isInStock bool) ([]models.Glasses, error) {
	return s.glassesRepo.GetGlassesByStock(ctx, page, pageSize, orderBy, sortBy, isInStock)
}

func (s *Service) GetSum() (int, error) {
	total, err := s.glassesRepo.GetSum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (s *Service) GetSumByType(glassesType string) (int, error) {
	total, err := s.glassesRepo.GetSumByType(context.Background(), glassesType)
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (s *Service) GetSumByStock(isInStock bool) (int, error) {
	total, err := s.glassesRepo.GetSumByStock(context.Background(), isInStock)
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (s *Service) GetGlassesReference(ctx context.Context, id uuid.UUID) (string, error) {
	return s.glassesRepo.GetGlassesReference(ctx, id)
}

package services

import (
	"context"
	"math"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/google/uuid"
)

func (s *Service) GetGlasses(ctx context.Context, page, pageSize int,
	orderBy, sortBy string) ([]models.Glasses, error) {
	return s.glassesRepo.GetGlasses(ctx, page, pageSize, orderBy, sortBy)
}

func (s *Service) GetGlassesByID(ctx context.Context, glassesID uuid.UUID) (*models.Glasses, error) {
	return s.glassesRepo.GetGlassesByID(ctx, glassesID)
}

func (s *Service) DeleteGlasses(ctx context.Context, glassesID uuid.UUID) error {
	return s.glassesRepo.DeleteGlasses(ctx, glassesID)
}

func (s *Service) UpdateGlasses(ctx context.Context, g models.Glasses) error {
	return s.glassesRepo.UpdateGlasses(ctx, g)
}

func (s *Service) InsertGlasses(ctx context.Context, g models.Glasses) error {
	return s.glassesRepo.InsertGlasses(ctx, g)
}

func (s *Service) GetSum() (int, error) {
	total, err := s.glassesRepo.GetSum(context.Background())
	pageSize := 20
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

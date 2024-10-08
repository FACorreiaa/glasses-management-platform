package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
)

func (s *Service) InsertShippingDetails(ctx context.Context, g, u uuid.UUID, c models.CustomerShippingForm,
	sh models.Shipping) error {
	return s.customerRepo.InsertShippingDetails(ctx, g, u, c, sh)
}

func (s *Service) GetCardIDNumber(ctx context.Context, u uuid.UUID) (string, error) {
	return s.customerRepo.GetCardIDNumber(ctx, u)
}

func (s *Service) GetShippingDetails(ctx context.Context, page, pageSize int,
	orderBy, sortBy, name string, leftEye, rightEye *float64) ([]models.ShippingDetails, error) {
	return s.customerRepo.GetShippingDetails(ctx, page, pageSize, orderBy, sortBy, name, leftEye, rightEye)
}

func (s *Service) GetShippingExpandedDetails(ctx context.Context, page, pageSize int,
	orderBy, sortBy, name string, leftEye, rightEye *float64) ([]models.SettingsShippingDetails, error) {
	return s.customerRepo.GetShippingExpandedDetails(ctx, page, pageSize, orderBy, sortBy, name, leftEye, rightEye)
}

func (s *Service) DeleteCustomer(ctx context.Context, customerID uuid.UUID) error {
	return s.glassesRepo.DeleteCustomer(ctx, customerID)
}

func (s *Service) UpdateShippingDetails(ctx context.Context, g models.ShippingDetailsForm, id uuid.UUID) error {
	return s.customerRepo.UpdateShippingDetails(ctx, g, id)
}

func (s *Service) GetCustomerGlassesID(ctx context.Context, customerID uuid.UUID) (*models.ShippingDetails, error) {
	return s.customerRepo.GetCustomerGlassesID(ctx, customerID)
}

func (s *Service) GetCardIDFromShipping(ctx context.Context, customerID uuid.UUID) (string, error) {
	return s.customerRepo.GetCardIDFromShipping(ctx, customerID)
}

func (s *Service) GetReferenceNumberFromShipping(ctx context.Context, customerID uuid.UUID) (string, error) {
	return s.customerRepo.GetReferenceNumberFromShipping(ctx, customerID)
}

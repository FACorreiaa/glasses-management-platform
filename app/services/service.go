package services

import (
	"github.com/FACorreiaa/glasses-management-platform/app/repository"
)

type Service struct {
	accountRepo  *repository.AccountRepository
	glassesRepo  *repository.GlassesRepository
	adminRepo    *repository.AdminRepository
	customerRepo *repository.CustomerRepository
}

func NewService(
	accountRepo *repository.AccountRepository,
	glassesRepo *repository.GlassesRepository,
	adminRepo *repository.AdminRepository,
	customerRepo *repository.CustomerRepository) *Service {
	return &Service{
		accountRepo:  accountRepo,
		glassesRepo:  glassesRepo,
		adminRepo:    adminRepo,
		customerRepo: customerRepo,
	}
}

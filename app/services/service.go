package services

import (
	"log"

	"github.com/FACorreiaa/glasses-management-platform/app/repository"
)

type Service struct {
	accountRepo *repository.AccountRepository
	glassesRepo *repository.GlassesRepository
}

func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func NewService(
	accountRepo *repository.AccountRepository,
	glassesRepo *repository.GlassesRepository) *Service {

	return &Service{
		accountRepo: accountRepo,
		glassesRepo: glassesRepo,
	}
}

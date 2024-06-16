package services

import (
	"log"

	"github.com/FACorreiaa/Aviation-tracker/app/repository"
)

type Service struct {
	accountRepo *repository.AccountRepository
}

func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func NewService(

	accountRepo *repository.AccountRepository) *Service {

	return &Service{

		accountRepo: accountRepo,
	}
}

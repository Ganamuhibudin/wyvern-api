package services

import (
	"net/http"
	"wyvern-api/models"
	"wyvern-api/utils"
)

// CarProcessor interface
type CarProcessor interface {
	FindAll() ([]models.Car, error)
}

// CarService struct
type CarService struct {
	identifier int64
	cp         CarProcessor
}

// NewCarService initiate CarService
func NewCarService(cp CarProcessor, identifier int64) *CarService {
	return &CarService{
		identifier: identifier,
		cp:         cp,
	}
}

// Find is method for find cars
func (svc *CarService) Find() models.Response {
	log := utils.NewLoggerIdentifier("Find", 0, svc.identifier).Service()
	cars, err := svc.cp.FindAll()
	if err != nil {
		log.Warn("failed get cars, error: %s", err.Error())
		return models.ResponseError(http.StatusNotFound, "user not found")
	}

	log.Info("success")
	return models.ResponseSuccessWithData(cars)
}

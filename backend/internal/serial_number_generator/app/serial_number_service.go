package app

import (
	"github.com/mjannello/smartcompost-webapp/backend/internal/serial_number_generator"
)

type SerialNumberGeneratorService interface {
	New() (string, error)
}

type serialNumberGeneratorService struct {
	generator serial_number_generator.SerialNumberGenerator
}

func NewSerialNumberGeneratorService(generator serial_number_generator.SerialNumberGenerator) SerialNumberGeneratorService {
	return &serialNumberGeneratorService{
		generator: generator,
	}
}

func (s *serialNumberGeneratorService) New() (string, error) {
	return s.generator.Generate()
}

package generators

import (
	"github.com/google/uuid"
	"github.com/mjannello/smartcompost-webapp/backend/internal/serial_number_generator"
)

type uuidGenerator struct{}

func NewUUIDGenerator() serial_number_generator.SerialNumberGenerator {
	return &uuidGenerator{}
}

func (g *uuidGenerator) Generate() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

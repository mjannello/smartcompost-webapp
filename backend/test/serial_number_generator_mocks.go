package test

import (
	"github.com/stretchr/testify/mock"
)

type SerialNumberGeneratorServiceMock struct {
	mock.Mock
}

func (sng *SerialNumberGeneratorServiceMock) New() (string, error) {
	args := sng.Called()
	sn, _ := args.Get(0).(string)
	e, _ := args.Get(1).(error)
	return sn, e
}

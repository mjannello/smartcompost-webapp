package clock

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type ClockMock struct {
	mock.Mock
}

func (m *ClockMock) Time() time.Time {
	args := m.Called()
	t, _ := args.Get(0).(time.Time)
	return t
}

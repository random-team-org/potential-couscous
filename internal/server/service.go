package server

import (
	"devices/internal/types"
	"errors"
)

type Service struct {
}

func NewService(dbopts, memcached string) *Service {
	return &Service{}
}
func (s *Service) GetDevices() ([]types.Device, error) {
	return nil, errors.New("not implemented yet")
}

func (s *Service) AllocateDevice(deviceID, userID string) error {
	return errors.New("not implemented yet")
}

func (s *Service) FreeDevice(deviceID string) error {
	return errors.New("not implemented yet")
}

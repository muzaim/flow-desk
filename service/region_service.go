package service

import (
	"flow-desk/client"
	"flow-desk/dto"
)

type RegionService interface {
	GetProvinces() ([]dto.Province, error)
}

type regionService struct {
	wilayahClient client.WilayahClient
}

func NewRegionService(wilayahClient client.WilayahClient) RegionService {
	return &regionService{wilayahClient: wilayahClient}
}

func (s *regionService) GetProvinces() ([]dto.Province, error) {
	provinces, err := s.wilayahClient.FetchProvinces()
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

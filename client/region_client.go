package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"flow-desk/dto"
)

type WilayahClient interface {
	FetchProvinces() ([]dto.Province, error)
}

type wilayahClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewWilayahClient() WilayahClient {
	return &wilayahClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://wilayah.id/api",
	}
}

func (c *wilayahClient) FetchProvinces() ([]dto.Province, error) {
	url := c.baseURL + "/provinces.json"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("gagal menghubungi server wilayah.id")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("terjadi kesalahan pada API wilayah.id")
	}

	var apiResponse dto.WilayahAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, errors.New("gagal membaca response dari wilayah.id")
	}

	return apiResponse.Data, nil
}

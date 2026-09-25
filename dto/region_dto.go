package dto

type Province struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type WilayahAPIResponse struct {
	Data []Province `json:"data"`
}

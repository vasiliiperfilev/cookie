package data

import "github.com/vasiliiperfilev/cookie/internal/validator"

type Item struct {
	Id         int64   `json:"id"`
	SupplierId int64   `json:"supplierId"`
	Unit       string  `json:"unit"`
	Size       float32 `json:"size"`
	Name       string  `json:"name"`
	ImageId    string  `json:"imageId"`
}

type PostItemDto struct {
	Unit    string  `json:"unit"`
	Size    float32 `json:"size"`
	Name    string  `json:"name"`
	ImageId string  `json:"imageId"`
}

func ValidatePostItemInput(v *validator.Validator, input PostItemDto) {
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(input.Unit != "", "unit", "must be provided")
	v.Check(input.Size > 0, "size", "must be positive number")
}

package models

type Company struct {
	ID                string `json:"id"`
	Name              string `json:"name" validate:"required,max=15"`
	Description       string `json:"description" validate:"max=3000"`
	AmountOfEmployees *int64 `json:"amountOfEmployees" validate:"required"`
	Registered        bool   `json:"registered"`
	CompanyType       string `json:"companyType" validate:"required,company-types"`
}

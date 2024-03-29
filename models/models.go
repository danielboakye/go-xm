package models

type Company struct {
	ID                string `json:"id"`
	Name              string `json:"name" validate:"required,max=15"`
	Description       string `json:"description" validate:"max=3000"`
	AmountOfEmployees int64  `json:"amountOfEmployees" validate:"gte=0"`
	Registered        bool   `json:"registered"`
	CompanyType       string `json:"companyType" validate:"required,eq=Corporations|eq=Non Profit|eq=Cooperative|eq=Sole Proprietorship"`
}

type CompanyUpdateReq struct {
	Name              *string `json:"name" validate:"omitempty,max=15"`
	Description       *string `json:"description" validate:"omitempty,max=3000"`
	AmountOfEmployees *int64  `json:"amountOfEmployees" validate:"omitempty,gte=0"`
	Registered        *bool   `json:"registered"`
	CompanyType       *string `json:"companyType" validate:"omitempty,eq=Corporations|eq=Non Profit|eq=Cooperative|eq=Sole Proprietorship"`
}

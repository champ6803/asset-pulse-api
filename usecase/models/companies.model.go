package models

type GetCompaniesResponse struct {
	Companies []CompanyItem `json:"companies"`
}

type CompanyItem struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

package import_gateway

import "time"

type ImportUsersData struct {
	UserName       string `json:"userName" validate:"required"`
	Email          string `json:"email" validate:"required, email"`
	DepartmentName string `json:"departmentName" validate:"required"`
}

type ImportUsersInput struct {
	CompanyID string     `json:"companyId" validate:"required"`
	Schedule  *time.Time `json:"schedule"`
	Data      []ImportUsersData
}

type ImportUsersOutput struct {
	ID           string
	DepartmentID string
	Name         string
	Email        string
	Login        string
}

type IImportGateway interface {
	ImportUsers(input *ImportUsersInput) ([]ImportUsersOutput, error)
}

package department_gateway

import (
	"conformity-core/enums/user_company_enum"
	"time"
)

type UserDepartment struct {
	Login           string
	Name            string
	Password        string
	UserID          string
	OID             string
	DepartmentID    string
	CompanyID       string
	RequirePassword bool
	SSOEnabled      bool
	Workspace       string
	Status          user_company_enum.Status
	Role            user_company_enum.Role
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

type CreateDepartmentGatewayInput struct {
	Name      string
	CompanyID string
}

type GetByIdDepartmentGatewayOutput struct {
	ID        string
	Name      string
	CompanyID string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type GetUsersByFilterGatewayInput struct {
	CompanyID         string
	Departments       []string
	UserCreatedAtGTE  time.Time
	UserCreatedAtLTE  time.Time
	UserCompanyStatus user_company_enum.Status
}

type GetUsersByFilterGatewayOutput struct {
	UserID       string
	DepartmentID string
	Name         string
	Email        string
}

type ListInput struct {
	CompanyID string
}

type ListOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeleteInput struct {
	CompanyID    string
	DepartmentID string
}

type PatchInput struct {
	ID        string
	CompanyID string
	Name      string
}

type GetOneByFilterInput struct {
	Name      string
	CompanyID string
}

type GetOneByFilterOutput struct {
	ID   string
	Name string
}

type PaginateUsersInput struct {
	CompanyID    string
	Name         string
	DepartmentID string
	Status       string
	Role         user_company_enum.Role
	Email        string

	ActualPage int
	PageSize   int
}

type PaginateUsers struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	Email          string                   `json:"email"`
	Phone          string                   `json:"phone"`
	JobPosition    string                   `json:"jobPosition"`
	Document       string                   `json:"document"`
	Status         user_company_enum.Status `json:"status"`
	Role           user_company_enum.Role   `json:"role"`
	DepartmentID   string                   `json:"departmentId"`
	DepartmentName string                   `json:"departmentName"`
	CreatedAt      time.Time                `json:"createdAt"`
}

type PaginateUsersOutput struct {
	Data     []PaginateUsers `json:"data"`
	MaxPages int             `json:"maxPages"`
}

type CountMonthlyUsersOutput struct {
	Count       int
	CompanyID   string
	CompanyName string
}

type CountMonthlyUsersInput struct {
	StartDateGTE time.Time
	EndDateLTE   time.Time
}

type IDepartmentGateway interface {
	GetUserByLogin(login string) (*UserDepartment, error)
	Create(input CreateDepartmentGatewayInput) (string, error)
	Delete(input DeleteInput) (bool, error)
	GetById(id string) (*GetByIdDepartmentGatewayOutput, error)
	PaginateUsers(input PaginateUsersInput) (*PaginateUsersOutput, error)
	GetOneByFilter(input GetOneByFilterInput) (*GetOneByFilterOutput, error)
	GetUsersByFilter(input GetUsersByFilterGatewayInput) (*[]GetUsersByFilterGatewayOutput, error)
	List(input ListInput) ([]ListOutput, error)
	Patch(input PatchInput) (bool, error)
	CountMonthlyUsers(input CountMonthlyUsersInput) ([]CountMonthlyUsersOutput, error)
}

package user_gateway

import (
	"conformity-core/enums/user_company_enum"
	"time"
)

type CreateInput struct {
	OID          string
	Name         string
	Password     string
	Email        string
	Login        string
	DepartmentID string
	Document     string
	Phone        string
	BirthDate    string
	JobPosition  string
	Role         user_company_enum.Role
}

type GetByIdUserGatewayOutput struct {
	ID          string
	Name        string
	Password    string
	Email       string
	Phone       string
	BirthDate   string
	Document    string
	JobPosition string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type GetByEmailUserGatewayOutput struct {
	ID          string
	Name        string
	Password    string
	Email       string
	Phone       string
	BirthDate   string
	Document    string
	JobPosition string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type ChangePasswordInput struct {
	Password string
	UserID   string
}

type PatchInput struct {
	ID           string
	OID          string
	DepartmentID string
	Name         string
	Document     string
	Email        string
	Login        string
	Phone        string
	JobPosition  string
	Role         user_company_enum.Role
}

type GetOneByFilterInput struct {
	ID           string
	DepartmentID string
}

type GetOneByFilterOutput struct {
	Phone string
}

type IUserGateway interface {
	Create(input CreateInput) (string, error)
	Patch(input PatchInput) (bool, error)
	GetById(id string) (*GetByIdUserGatewayOutput, error)
	GetOneByFilter(id GetOneByFilterInput) (*GetOneByFilterOutput, error)
	GetByEmail(email string) (*GetByEmailUserGatewayOutput, error)
	ChangePassword(input ChangePasswordInput) (bool, error)
}

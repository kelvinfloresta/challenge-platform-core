package company_gateway

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	"time"

	"gorm.io/gorm"
)

type CreateCompanyGatewayInput struct {
	Name            string
	Document        string
	Workspace       string
	Domain          string
	RequirePassword bool
}

type GetByIdOutput struct {
	ID              string
	Name            string
	Document        string
	Workspace       string
	IDPMetadata     string
	RequirePassword bool
	SSOEnabled      bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type GetUserInput struct {
	DepartmentID string
	UserID       string
}

type CompanyUser struct {
	UserID       string
	DepartmentID string
	Status       user_company_enum.Status
	Role         user_company_enum.Role
}

type ChangeUserRoleInput struct {
	DepartmentID string
	UserID       string
	Role         user_company_enum.Role
}

type GetUserOutput struct {
	UserID       string
	DepartmentID string
	Status       user_company_enum.Status
	Role         user_company_enum.Role
}

type ListManagersInput struct {
	CompanyIDs []string
	Status     user_company_enum.Status
}

type ListManagersOutput struct {
	CompanyID string
	Email     string
}

type GetOneByFilterInput struct {
	DepartmentID string
	Workspace    string
	Domain       string
}

type GetOneByFilterOutput struct {
	ID                       string
	Name                     string
	RequirePassword          bool
	SSOEnabled               bool
	Workspace                string
	IdentityProviderMetadata string
}

type ICompanyGateway interface {
	GetTransaction(ctx *context.CoreCtx) *gorm.DB
	Create(input CreateCompanyGatewayInput) (string, error)
	GetById(id string) (*GetByIdOutput, error)
	GetUser(input GetUserInput) (*GetUserOutput, error)
	ListManagers(input ListManagersInput) ([]ListManagersOutput, error)
	ChangeUserRole(input ChangeUserRoleInput, tx *gorm.DB) (bool, error)
	GetOneByFilter(input GetOneByFilterInput) (*GetOneByFilterOutput, error)
}

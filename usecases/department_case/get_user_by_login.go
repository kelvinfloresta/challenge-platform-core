package department_case

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
}

func (d DepartmentCase) GetUserByLogin(login string) (*UserDepartment, error) {
	user, err := d.gateway.GetUserByLogin(login)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &UserDepartment{
		Login:           user.Login,
		Name:            user.Name,
		Password:        user.Password,
		CompanyID:       user.CompanyID,
		UserID:          user.UserID,
		DepartmentID:    user.DepartmentID,
		SSOEnabled:      user.SSOEnabled,
		Workspace:       user.Workspace,
		Status:          user.Status,
		Role:            user.Role,
		OID:             user.OID,
		RequirePassword: user.RequirePassword,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}, nil
}

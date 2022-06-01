package department_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	g "conformity-core/gateways/department_gateway"
)

type PaginateUsersInput struct {
	ActualPage   int `json:"actual_page"`
	PageSize     int `json:"page_size"`
	Name         string
	DepartmentID string
	Status       string
	Role         user_company_enum.Role
	Email        string
}

func (d DepartmentCase) PaginateUsers(ctx *context.CoreCtx, input PaginateUsersInput) (*g.PaginateUsersOutput, error) {
	return d.gateway.PaginateUsers(g.PaginateUsersInput{
		CompanyID:    ctx.Session.CompanyID,
		Name:         input.Name,
		DepartmentID: input.DepartmentID,
		Status:       input.Status,
		Role:         input.Role,
		Email:        input.Email,
		ActualPage:   input.ActualPage,
		PageSize:     input.PageSize,
	})
}

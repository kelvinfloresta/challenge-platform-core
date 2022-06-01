package department_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	g "conformity-core/gateways/department_gateway"

	"github.com/sirupsen/logrus"
)

func (d DepartmentCase) Delete(ctx *context.CoreCtx, departmentId string) (bool, error) {
	if ctx.Session.Role != user_company_enum.RoleCompanyManager {
		return false, core_errors.ErrWithoutPrivileges
	}

	users, err := d.GetUsersByFilter(GetUsersByFilterCaseInput{
		CompanyID:   ctx.Session.CompanyID,
		Departments: []string{departmentId},
	})

	if err != nil {
		return false, err
	}

	if users != nil && len(*users) > 0 {
		return false, core_errors.ErrDeleteDepartmentWithUsers
	}

	deleted, err := d.gateway.Delete(g.DeleteInput{
		CompanyID:    ctx.Session.CompanyID,
		DepartmentID: departmentId,
	})

	if err != nil {
		return false, err
	}

	ctx.Logger.WithFields(logrus.Fields{
		"id":      departmentId,
		"deleted": deleted,
	}).Info("delete-department")

	return deleted, nil
}

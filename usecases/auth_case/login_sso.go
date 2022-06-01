package auth_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	"conformity-core/usecases/company_case"
	"conformity-core/usecases/department_case"
	"conformity-core/usecases/user_case"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
)

type LoginSSOInput struct {
	OID            string
	Email          string
	Name           string
	DepartmentName string
	Workspace      string
}

func (a *AuthCase) LoginSSO(input LoginSSOInput) (string, error) {
	if input.Email == "" {
		return "", errors.New("MISSING_EMAIL")
	}

	user, err := a.departmentCase.GetUserByLogin(input.Email)

	if err != nil {
		return "", err
	}

	if user == nil {
		user, err = a.createSSOUser(input)
		if err != nil {
			return "", err
		}
	}

	if user.Status == user_company_enum.Suspended {
		return "", ErrSuspendedUser
	}

	shouldResetOID := user.OID == ""
	if shouldResetOID {
		err := a.resetOID(input.OID, user.UserID, user.DepartmentID)
		if err != nil {
			return "", err
		}
		return generateToken(user, nil)
	}

	if user.OID != input.OID {
		return "", nil
	}

	return generateToken(user, nil)
}

func (a AuthCase) createSSOUser(input LoginSSOInput) (out *department_case.UserDepartment, err error) {
	if input.DepartmentName == "" {
		return nil, core_errors.ErrUserWithoutDepartment
	}

	company, err := a.companyCase.GetOneByFilter(company_case.GetOneByFilterInput{
		Workspace: input.Workspace,
	})

	if err != nil {
		return
	}

	if company == nil {
		sentry.CaptureMessage("Workspace not found: " + input.Workspace)
		return nil, core_errors.ErrCompanyNotFound
	}

	departmentId, err := a.createOrGetDepartment(company.ID, input.DepartmentName)

	if err != nil {
		return
	}

	defaultRole := user_company_enum.RoleUser
	login := input.Email
	userId, err := a.userCase.Create(context.Internal, user_case.CreateInput{
		OID:          input.OID,
		Name:         input.Name,
		Email:        input.Email,
		Login:        login,
		DepartmentID: departmentId,
		Role:         defaultRole,
	})

	if err != nil {
		return
	}

	out = &department_case.UserDepartment{
		OID:          input.OID,
		UserID:       userId,
		Login:        login,
		Name:         input.Name,
		DepartmentID: departmentId,
		CompanyID:    company.ID,
		Role:         defaultRole,
	}

	return
}

func (a AuthCase) createOrGetDepartment(companyId, departmentName string) (string, error) {
	department, err := a.departmentCase.GetOneByFilter(department_case.GetOneByFilter{
		Name:      departmentName,
		CompanyID: companyId,
	})

	if err != nil {
		return "", err
	}

	if department != nil {
		return department.ID, nil
	}

	return a.departmentCase.Create(context.Internal, department_case.CreateDepartmentCaseInput{
		Name:      departmentName,
		CompanyID: companyId,
	})
}

func (a AuthCase) resetOID(newOID, userId, departmentId string) error {
	updated, err := a.userCase.Patch(context.Internal, user_case.PatchInput{
		ID:           userId,
		DepartmentID: departmentId,

		OID: newOID,
	})

	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	if !updated {
		message := fmt.Sprintf("Cannot update OID of User:%s Department:%s", userId, departmentId)
		sentry.CaptureMessage(message)
	}

	return nil
}

package campaign_case

import (
	"conformity-core/context"
	"conformity-core/enums/campaign_enum"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	"conformity-core/gateways/campaign_gateway"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RemoveUserInput struct {
	UserID       string `json:"userId" validate:"required"`
	DepartmentID string `json:"departmentId" validate:"required"`
}

func (c CampaignCase) RemoveUser(
	ctx *context.CoreCtx,
	input RemoveUserInput,
) (bool, error) {
	isValid, err := c.validatePermissions(ctx, input.DepartmentID)
	if err != nil {
		return false, err
	}

	if !isValid {
		return false, nil
	}

	removed := false
	suspendedResults := false
	err = c.gateway.Transaction(func(tx *gorm.DB) error {
		removed, err = c.gateway.RemoveUser(ctx, campaign_gateway.RemoveUserInput{
			DepartmentID: input.DepartmentID,
			UserID:       input.UserID,
		}, tx)

		if err != nil {
			return err
		}

		if !removed {
			return nil
		}

		suspendedResults, err = c.UpdateResult(UpdateResultInput{
			DepartmentID: input.DepartmentID,
			UserID:       input.UserID,
			NewStatus:    campaign_enum.Suspended,
		}, tx)

		return err
	})

	if err != nil {
		return false, err
	}

	ctx.Logger.WithFields(logrus.Fields{
		"user-removed":            input.UserID,
		"user-removed-department": input.DepartmentID,
		"removed":                 removed,
		"suspendedResults":        suspendedResults,
	}).Info("remove-user")

	return removed, nil
}

func (c CampaignCase) validatePermissions(ctx *context.CoreCtx, departmentId string) (bool, error) {
	if ctx.Session.Role != user_company_enum.RoleBackoffice &&
		ctx.Session.Role != user_company_enum.RoleCompanyManager {
		return false, core_errors.ErrWithoutPrivileges
	}

	if ctx.Session.Role == user_company_enum.RoleBackoffice {
		return true, nil
	}

	company, err := c.departmentCase.GetById(departmentId)

	if err != nil {
		return false, err
	}

	isValid := company != nil && company.CompanyID == ctx.Session.CompanyID
	return isValid, nil
}

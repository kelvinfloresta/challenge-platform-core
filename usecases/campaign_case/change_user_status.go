package campaign_case

import (
	"conformity-core/context"
	"conformity-core/enums/campaign_enum"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	g "conformity-core/gateways/campaign_gateway"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ChangeUserStatusCaseInput struct {
	DepartmentID string                   `json:"department_id" validate:"required"`
	UserID       string                   `json:"user_id" validate:"required"`
	Status       user_company_enum.Status `json:"status" validate:"required"`
}

func (c CampaignCase) ChangeUserStatus(ctx *context.CoreCtx, input ChangeUserStatusCaseInput) (updated bool, err error) {
	if ctx.Session.Role != user_company_enum.RoleCompanyManager &&
		ctx.Session.Role != user_company_enum.RoleBackoffice {
		return false, core_errors.ErrWithoutPrivileges
	}

	updatedResult := false
	err = c.gateway.Transaction(func(tx *gorm.DB) error {
		updated, err = c.gateway.ChangeUserStatus(g.ChangeUserStatusInput{
			UserID:       input.UserID,
			DepartmentID: input.DepartmentID,
			Status:       input.Status,
		}, tx)

		if err != nil {
			return err
		}

		if !updated {
			return nil
		}

		updatedResult, err = c.UpdateResult(UpdateResultInput{
			UserID:       input.UserID,
			DepartmentID: input.DepartmentID,
			NewStatus:    parseResultStatus(input.Status),
		}, tx)

		return err
	})

	if err != nil {
		return false, err
	}

	ctx.Logger.WithFields(logrus.Fields{
		"status":        input.Status,
		"updated":       updated,
		"updatedResult": updatedResult,
	}).Info("change-user-status")

	return updated, err
}

func parseResultStatus(status user_company_enum.Status) campaign_enum.ResultStatus {
	switch status {
	case user_company_enum.Active:
		return campaign_enum.Active

	case user_company_enum.Suspended:
		return campaign_enum.Suspended

	default:
		return campaign_enum.Active
	}
}

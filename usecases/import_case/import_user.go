package import_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	e "conformity-core/errors"
	g "conformity-core/gateways/import_gateway"
	"conformity-core/usecases/user_case"
	"time"

	"github.com/sirupsen/logrus"
)

func (i *ImportCase) ImportUsers(ctx *context.CoreCtx, input *g.ImportUsersInput) error {
	if ctx.Session.Role != user_company_enum.RoleBackoffice {
		return e.ErrWithoutPrivileges
	}

	imported, err := i.gateway.ImportUsers(input)
	if err != nil {
		return err
	}

	go i.notifyImportedUsers(imported, input.Schedule)

	go i.logImportedUsers(ctx.Logger, input.CompanyID, imported, input.Schedule)

	return nil
}

func (i *ImportCase) logImportedUsers(
	logger *logrus.Entry,
	companyId string,
	imported []g.ImportUsersOutput,
	schedule *time.Time,
) {
	logger.WithFields(logrus.Fields{
		"total":     len(imported),
		"schedule":  schedule,
		"companyId": companyId,
	}).Info("import-users-resume")

	for _, user := range imported {
		logger.WithFields(logrus.Fields{
			"user-id":            user.ID,
			"user-department-id": user.DepartmentID,
			"name":               user.Name,
			"email":              user.Email,
			"login":              user.Login,
		}).Info("import-users")
	}
}

func (i *ImportCase) notifyImportedUsers(imported []g.ImportUsersOutput, schedule *time.Time) {
	notifications := make([]*user_case.NotifyNewAccountInput, 0, len(imported))
	for _, iuo := range imported {
		notifications = append(notifications, &user_case.NotifyNewAccountInput{
			UserID:       iuo.ID,
			DepartmentID: iuo.DepartmentID,
			Email:        iuo.Email,
			Login:        iuo.Login,
			Name:         iuo.Name,
			Schedule:     schedule,
		})
	}

	i.userCase.BatchNotifyNewAccount(notifications)
}

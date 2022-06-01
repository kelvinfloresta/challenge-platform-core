package department_case

import (
	"conformity-core/context"
	"conformity-core/gateways/department_gateway"

	"github.com/sirupsen/logrus"
)

type PatchInput struct {
	ID   string
	Name string
}

func (d DepartmentCase) Patch(ctx *context.CoreCtx, input PatchInput) (bool, error) {
	updated, err := d.gateway.Patch(department_gateway.PatchInput{
		ID:        input.ID,
		CompanyID: ctx.Session.CompanyID,
		Name:      input.Name,
	})

	if err != nil {
		return false, err
	}

	ctx.Logger.WithFields(logrus.Fields{
		"id":        input.ID,
		"companyId": ctx.Session.CompanyID,
		"name":      input.Name,
		"updated":   updated,
	})

	return updated, nil
}

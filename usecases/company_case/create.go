package company_case

import (
	"conformity-core/context"
	g "conformity-core/gateways/company_gateway"

	"github.com/sirupsen/logrus"
)

type CreateCompanyCaseInput struct {
	Name            string `json:"name" validate:"required"`
	Document        string `json:"document" validate:"required"`
	Workspace       string `json:"workspace" validate:"required"`
	Domain          string `json:"domain"`
	RequirePassword bool   `json:"requirePassword"`
}

func (u CompanyCase) Create(ctx *context.CoreCtx, input CreateCompanyCaseInput) (string, error) {
	id, err := u.gateway.Create(g.CreateCompanyGatewayInput(input))
	if err != nil {
		return "", err
	}

	ctx.Logger.WithFields(logrus.Fields{
		"name":            input.Name,
		"document":        input.Document,
		"workspace":       input.Workspace,
		"domain":          input.Domain,
		"requirePassword": input.RequirePassword,
	}).Info("create-company")

	return id, nil
}

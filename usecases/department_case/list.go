package department_case

import (
	"conformity-core/context"
	g "conformity-core/gateways/department_gateway"
)

func (u DepartmentCase) List(ctx *context.CoreCtx) ([]g.ListOutput, error) {
	return u.gateway.List(g.ListInput{
		CompanyID: ctx.Session.CompanyID,
	})
}

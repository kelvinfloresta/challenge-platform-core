package department_case

import (
	g "conformity-core/gateways/department_gateway"
	"time"

	now_lib "github.com/jinzhu/now"
)

func (d DepartmentCase) CountMonthlyUsers(now time.Time) ([]g.CountMonthlyUsersOutput, error) {
	nt := now_lib.With(now)
	return d.gateway.CountMonthlyUsers(g.CountMonthlyUsersInput{
		StartDateGTE: nt.BeginningOfMonth(),
		EndDateLTE:   nt.EndOfMonth(),
	})
}

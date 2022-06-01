package user_case

import (
	"conformity-core/usecases/company_case"
	"fmt"

	"github.com/getsentry/sentry-go"
)

func (u *UserCase) BatchNotifyNewAccount(input []*NotifyNewAccountInput) {
	hub := sentry.CurrentHub().Clone()

	if len(input) == 0 {
		hub.CaptureMessage("BatchNotifyNewAccountInput with empty slice")
		return
	}

	departmentId := input[0].DepartmentID
	company, err := u.companyCase.GetOneByFilter(company_case.GetOneByFilterInput{
		DepartmentID: departmentId,
	})

	if err != nil {
		hub.CaptureException(fmt.Errorf("GetOneByFilter: %v", err))
		return
	}

	if company == nil {
		hub.CaptureMessage("Company not found with DepartmentID: " + departmentId)
		return
	}

	if company.RequirePassword {
		for _, el := range input {
			go u.notifyNewAccount(el, hub)
		}
		return
	}

	for _, el := range input {
		go u.notifyNewAccountWithoutPassword(el, hub)
	}
}

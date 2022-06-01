package billing_case

import (
	"conformity-core/context"
	"fmt"
	"time"

	now_lib "github.com/jinzhu/now"
)

func (b BillingCase) SendBillingReport(now time.Time) error {
	lastMonth := now_lib.With(now).BeginningOfMonth().AddDate(0, 0, -1)

	result, err := b.departmentCase.CountMonthlyUsers(lastMonth)
	if err != nil {
		return fmt.Errorf("CountMonthlyUsers: %v", err)
	}

	err = b.notificationCase.BillingReport(context.Internal, result)
	if err != nil {
		return fmt.Errorf("BillingReport: %v", err)
	}

	return nil
}

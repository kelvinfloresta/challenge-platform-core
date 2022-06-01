package billing_case

import (
	"conformity-core/usecases/department_case"
	"conformity-core/usecases/notification_case"
)

type BillingCase struct {
	notificationCase *notification_case.NotificationCase
	departmentCase   *department_case.DepartmentCase
}

func New(
	notificationCase *notification_case.NotificationCase,
	departmentCase *department_case.DepartmentCase,
) BillingCase {
	return BillingCase{notificationCase, departmentCase}
}

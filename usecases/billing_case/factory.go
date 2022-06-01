package billing_case

import (
	"conformity-core/usecases/department_case"
	"conformity-core/usecases/notification_case"
)

var Singleton = build()

func build() BillingCase {
	return New(notification_case.Singleton, &department_case.Singleton)
}

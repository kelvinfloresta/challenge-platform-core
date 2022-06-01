package notification_case

import (
	"bytes"
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/gateways/department_gateway"
	"conformity-core/gateways/notification_gateway"
	"encoding/csv"
	"fmt"
)

func billingReportToCSV(users []department_gateway.CountMonthlyUsersOutput) ([]byte, error) {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	if err := w.Error(); err != nil {
		return []byte{}, err
	}

	if err := w.Write([]string{"ID", "Empresa", "Total"}); err != nil {
		return []byte{}, err
	}

	for _, u := range users {
		if err := w.Write([]string{u.CompanyID, u.CompanyName, fmt.Sprint(u.Count)}); err != nil {
			return []byte{}, err
		}
	}

	w.Flush()

	return b.Bytes(), nil
}

func (n *NotificationCase) BillingReport(
	ctx *context.CoreCtx,
	input []department_gateway.CountMonthlyUsersOutput,
) error {
	csv, err := billingReportToCSV(input)
	if err != nil {
		return err
	}

	return n.sendEmail(ctx, notification_gateway.SendEmailInput{
		To:         []string{config.BillingEmail},
		Subject:    "Billing Report",
		TemplateID: "billing-report",
		Attachment: &notification_gateway.Attachment{
			Name: "report.csv",
			File: csv,
		},
	})
}

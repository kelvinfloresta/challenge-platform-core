package notification_case

import (
	"bytes"
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/frameworks/logger"
	"conformity-core/gateways/notification_gateway"
	"encoding/csv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

type DeadlineChallengeUser struct {
	Name           string
	Email          string
	DepartmentName string
}

type ManagerDeadlineChallengeInput struct {
	CompanyID  string
	Users      []DeadlineChallengeUser
	Managers   []string
	ReportDate time.Time
	Schedule   time.Time
}

func (n *NotificationCase) ManagerDeadlineChallenge(
	ctx *context.CoreCtx,
	input ManagerDeadlineChallengeInput,
) {
	const notificationName = "notify-manager-deadline-challenge"
	log := logger.Internal.WithFields(logrus.Fields{
		"notification": notificationName,
		"companyId":    input.CompanyID,
	})

	csv, err := listUsersToCSV(input.Users)

	if err != nil {
		log.Error("create-csv")
		sentry.CaptureException(err)
		return
	}

	date := input.ReportDate.Format("02/01")
	subject := "Conformity Pro - Relatório " + date
	err = n.sendEmail(ctx, notification_gateway.SendEmailInput{
		To:         input.Managers,
		Subject:    subject,
		TemplateID: "users-without-completed-challenge-report",
		Schedule:   &input.Schedule,
		Attachment: &notification_gateway.Attachment{
			File: csv,
			Name: "attachment.csv",
		},
		Variables: map[string]string{
			"link": config.AppURL,
		},
	})

	if err != nil {
		log.Error("send-email")
		sentry.CaptureException(err)
		return
	}

	log.Info("email-sended")

}

func listUsersToCSV(users []DeadlineChallengeUser) ([]byte, error) {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	if err := w.Error(); err != nil {
		return []byte{}, err
	}

	if err := w.Write([]string{"Nome", "Email", "Departamento"}); err != nil {
		return []byte{}, err
	}

	for _, u := range users {
		if err := w.Write([]string{u.Name, u.Email, u.DepartmentName}); err != nil {
			return []byte{}, err
		}
	}

	w.Flush()

	return b.Bytes(), nil
}

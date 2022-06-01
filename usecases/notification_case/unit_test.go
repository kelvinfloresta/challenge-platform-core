package notification_case_test

import (
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/fixtures"
	"conformity-core/gateways/notification_gateway"
	"conformity-core/usecases/notification_case"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DeadlineChallenge_should_do_happy_path(t *testing.T) {
	messages := []notification_gateway.SendEmailInput{}
	sendEmailSpy := func(ctx *context.CoreCtx, input notification_gateway.SendEmailInput) error {
		messages = append(messages, input)
		return nil
	}

	sut := notification_case.New(sendEmailSpy)

	email := "any@mail.com"
	name := "Any name"

	sut.DeadlineChallenge(fixtures.DUMMY_CONTEXT, email, name, nil)

	require.Len(t, messages, 1)
	for _, message := range messages {
		require.Equal(t, "deadline-challenge", message.TemplateID)
		require.Equal(t, "Você possui um desafio pendente", message.Subject)
		require.Len(t, message.To, 1)
		require.Equal(t, email, message.To[0])
		require.Equal(t, name, message.Variables["name"])
		require.Equal(t, config.AppURL, message.Variables["link"])
	}
}

func Test_NotifyManagerDeadlineChallenge__Should_do_happy_path(t *testing.T) {

	sut, sendMailMock := fixtures.NotificationCase.NewUnit()

	now := time.Now()
	managers := []string{
		"manager1@mail.com",
		"manager2@mail.com",
	}
	users := []notification_case.DeadlineChallengeUser{
		{
			Name:           "any_name",
			Email:          "any_mail@mail.com",
			DepartmentName: "any_department",
		},
	}

	sut.ManagerDeadlineChallenge(fixtures.DUMMY_CONTEXT, notification_case.ManagerDeadlineChallengeInput{
		CompanyID:  "any_company_id",
		Users:      users,
		Managers:   managers,
		ReportDate: now,
	})

	require.Len(t, sendMailMock.Messages, 1)
	message := sendMailMock.Messages[0]

	require.Equal(t, "users-without-completed-challenge-report", message.TemplateID)

	date := now.Format("02/01")
	expectedSubject := "Conformity Pro - Relatório " + date
	require.Equal(t, expectedSubject, message.Subject)

	require.Len(t, message.To, len(managers))

	for i, emailManager := range managers {
		require.Equal(t, emailManager, message.To[i])
		require.Equal(t, config.AppURL, message.Variables["link"])
	}

	csv := string(message.Attachment.File)
	for _, user := range users {
		row := fmt.Sprintf("%s,%s,%s", user.Name, user.Email, user.DepartmentName)

		assert.Contains(t, csv, row)
	}
}

func Test_NotifyNewChallenge__should_do_happy_path(t *testing.T) {
	messages := []notification_gateway.SendEmailInput{}
	sendEmailSpy := func(ctx *context.CoreCtx, input notification_gateway.SendEmailInput) error {
		messages = append(messages, input)
		return nil
	}

	sut := notification_case.New(sendEmailSpy)

	email := "any@mail.com"
	name := "Any name"

	schedule := time.Now()

	sut.NewChallenge(fixtures.DUMMY_CONTEXT, notification_case.NewChallengeInput{
		Users: []notification_case.NewChallengeUser{
			{
				Email: email,
				Name:  name,
			},
		},
		Schedule: &schedule,
	})

	require.Len(t, messages, 1)
	for _, message := range messages {
		require.Equal(t, "new-challenge", message.TemplateID)
		require.Equal(t, "Novo desafio enviado para você", message.Subject)
		require.Len(t, message.To, 1)
		require.Equal(t, email, message.To[0])
		require.Equal(t, name, message.Variables["name"])
		require.Equal(t, config.AppURL, message.Variables["link"])
		require.Equal(t, &schedule, message.Schedule)
	}

}

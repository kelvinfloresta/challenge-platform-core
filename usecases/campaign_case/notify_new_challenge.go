package campaign_case

import (
	coreContext "conformity-core/context"
	"conformity-core/frameworks/logger"
	"conformity-core/usecases/notification_case"
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

const newChallenge = "notify-start-challenge"

func (n CampaignCase) NotifyNewChallenge(now time.Time) error {
	ctx := coreContext.New(context.Background())

	withoutMinutes := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(), 0, 0, 0, time.UTC)

	tomorrow := withoutMinutes.AddDate(0, 0, 1)

	users, err := n.ListUsers(
		ctx,
		ListUsersInput{
			StartDate: tomorrow,
		},
	)

	if err != nil {
		logger.Internal.WithFields(logrus.Fields{
			"notification": newChallenge,
		}).Error("list-users")
		sentry.CaptureException(err)
		return err
	}

	tomorrowAt12AM := time.Date(
		now.Year(),
		now.Month(),
		now.Day()+1,
		12, 0, 0, 0, time.UTC)

	usersToNotify := make([]notification_case.NewChallengeUser, len(users))
	for i, u := range users {
		usersToNotify[i] = notification_case.NewChallengeUser{
			Email: u.Email,
			Name:  u.Name,
		}
	}

	n.notificationCase.NewChallenge(ctx, notification_case.NewChallengeInput{
		Users:    usersToNotify,
		Schedule: &tomorrowAt12AM,
	})

	logger.Internal.WithFields(logrus.Fields{
		"notification": newChallenge,
	}).Info("task-done")

	return nil
}

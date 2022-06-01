package campaign_case

import (
	coreContext "conformity-core/context"
	"conformity-core/domain/challenge"
	"conformity-core/frameworks/logger"
	"conformity-core/utils"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func (n CampaignCase) NotifyDeadlineChallenge(now time.Time) error {
	now = now.UTC()
	const notificationName = "notify-deadline-challenge"

	startOfDay := utils.StartOfDay(now)
	gt := startOfDay.AddDate(0, 0, 1).Add(12 * time.Hour)
	lte := startOfDay.AddDate(0, 0, 2).Add(12 * time.Hour)

	notCorrect := false
	maxTries := challenge.MAX_TRIES
	users, err := n.ListUsers(
		coreContext.Internal,
		ListUsersInput{
			StartDateLTE: now,
			EndDateGT:    gt,
			EndDateLTE:   lte,
			Correct:      &notCorrect,
			TriesLT:      &maxTries,
		},
	)

	if err != nil {
		logger.Internal.WithFields(logrus.Fields{
			"notification": notificationName,
		}).Error("list-users")
		sentry.CaptureException(err)
		return err
	}

	schedule := time.Date(
		now.Year(),
		now.Month(),
		now.Day()+1,
		12, 0, 0, 0, time.UTC)

	for _, u := range users {
		n.notificationCase.DeadlineChallenge(coreContext.Internal, u.Email, u.Name, &schedule)
	}

	logger.Internal.WithFields(logrus.Fields{
		"notification": notificationName,
	}).Info("task-done")

	return nil
}

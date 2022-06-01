package campaign_case

import (
	coreContext "conformity-core/context"
	"conformity-core/domain/challenge"
	"conformity-core/frameworks/logger"
	"conformity-core/usecases/notification_case"
	"conformity-core/utils"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func (n CampaignCase) NotifyManagerDeadlineChallenge(now time.Time) error {
	now = now.UTC()
	const notificationName = "notify-manager-deadline-challenge"

	log := logger.Internal.WithFields(logrus.Fields{
		"notification": notificationName,
	})

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
		log.Error("list-users")
		sentry.CaptureException(err)
		return err
	}

	companies := groupByCompanyID(users)
	companyIds := make([]string, 0, len(companies))
	for k := range companies {
		companyIds = append(companyIds, k)
	}

	managers, err := n.getManagers(companyIds)
	if err != nil {
		log.Error("get-managers")
		sentry.CaptureException(err)
		return err
	}

	schedule := time.Date(
		now.Year(),
		now.Month(),
		now.Day()+1,
		12, 0, 0, 0, time.UTC)

	for companyId, deadlineUsers := range companies {
		n.notificationCase.ManagerDeadlineChallenge(coreContext.Internal, notification_case.ManagerDeadlineChallengeInput{
			CompanyID:  companyId,
			Users:      deadlineUsers,
			Managers:   managers[companyId],
			ReportDate: now,
			Schedule:   schedule,
		})
	}

	logger.Internal.WithFields(logrus.Fields{
		"notification": notificationName,
	}).Info("task-done")

	return nil
}

func groupByCompanyID(users []ListUsersOutput) map[string][]notification_case.DeadlineChallengeUser {
	companies := make(map[string][]notification_case.DeadlineChallengeUser)
	for _, u := range users {
		companies[u.CompanyID] = append(companies[u.CompanyID], notification_case.DeadlineChallengeUser{
			Name:           u.Name,
			Email:          u.Email,
			DepartmentName: u.DepartmentName,
		})
	}

	return companies
}

func (c CampaignCase) getManagers(companyIds []string) (map[string][]string, error) {
	output := make(map[string][]string)
	managers, err := c.companyCase.ListManagers(companyIds)
	if err != nil {
		return nil, err
	}

	for _, manager := range managers {
		output[manager.CompanyID] = append(output[manager.CompanyID], manager.Email)
	}

	return output, nil
}

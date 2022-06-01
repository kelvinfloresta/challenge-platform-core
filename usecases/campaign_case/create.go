package campaign_case

import (
	"conformity-core/context"
	"conformity-core/enums/campaign_enum"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	"conformity-core/frameworks/logger"
	g "conformity-core/gateways/campaign_gateway"
	"conformity-core/gateways/challenge_gateway"
	"conformity-core/gateways/department_gateway"
	"conformity-core/usecases/department_case"
	"conformity-core/usecases/notification_case"
	"conformity-core/utils"
	"errors"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

type CreateScheduledChallenge struct {
	ChallengeID string    `validate:"required"`
	StartDate   time.Time `validate:"required,datetime"`
	EndDate     time.Time `validate:"required,datetime"`
}

type CreateCampaignCaseInput struct {
	Title                 string                     `json:"title"`
	Challenges            []CreateScheduledChallenge `json:"challenges" validate:"required"`
	OnlyUsersCreatedAtGTE time.Time                  `json:"onlyUsersCreatedAtGte"`
	OnlyUsersCreatedAtLTE time.Time                  `json:"onlyUsersCreatedAtLte"`
	OnlyDepartments       []string                   `json:"onlyDepartments"`
	CompanyID             string                     `json:"companyId"`
}

func (c CampaignCase) Create(ctx *context.CoreCtx, input CreateCampaignCaseInput) (id string, err error) {
	if ctx.Session.Role != user_company_enum.RoleCompanyManager &&
		ctx.Session.Role != user_company_enum.RoleBackoffice {
		return "", core_errors.ErrWithoutPrivileges
	}

	if ctx.Session.Role == user_company_enum.RoleCompanyManager {
		input.CompanyID = ctx.Session.CompanyID
	}

	if input.CompanyID == "" {
		return "", errors.New("MISSING_COMPANY_ID")
	}

	err = validateScheduleChallenges(input.Challenges)
	if err != nil {
		return
	}

	campaignUsers, err := c.getCampaignUsers(
		input.CompanyID,
		input.OnlyDepartments,
		input.OnlyUsersCreatedAtGTE,
		input.OnlyUsersCreatedAtLTE,
	)

	if err != nil {
		return
	}

	questions, err := c.getQuestions(input.Challenges)
	if err != nil {
		return
	}

	campaignId, err := utils.UUID()
	if err != nil {
		return
	}

	participants := []g.CreateCampaignParticipant{}
	for _, user := range *campaignUsers {
		participants = append(participants, g.CreateCampaignParticipant{
			UserID:       user.UserID,
			DepartmentID: user.DepartmentID,
		})
	}

	campaignId, err = c.gateway.Create(g.CreateCampaignGatewayInput{
		CampaignID:          campaignId,
		Title:               input.Title,
		Results:             makeResults(campaignId, *questions, participants),
		ScheduledChallenges: makeScheduledChallenge(campaignId, input.Challenges),
	})

	if err != nil {
		return
	}

	go c.notifyTodayChallenge(campaignId, input.Challenges, *campaignUsers)

	ctx.Logger.WithFields(logrus.Fields{
		"id":                    id,
		"companyId":             input.CompanyID,
		"onlyDepartments":       strings.Join(input.OnlyDepartments, ","),
		"onlyUsersCreatedAtGTE": input.OnlyUsersCreatedAtGTE,
		"onlyUsersCreatedAtLTE": input.OnlyUsersCreatedAtLTE,
	}).Info("campaign-created")

	return campaignId, err
}

func (c CampaignCase) notifyTodayChallenge(
	campaignId string,
	challenges []CreateScheduledChallenge,
	users []department_gateway.GetUsersByFilterGatewayOutput,
) {
	defer func() {
		if r := recover(); r != nil {
			errorField := "campaign-new-challenge-notify-" + campaignId

			log := logger.Internal.WithFields(logrus.Fields{
				"campaignId": campaignId,
			})

			if err, ok := r.(error); ok {
				log.WithError(err).Error(errorField)
				sentry.CaptureException(err)
				return
			}

			log.Error("campaign-new-challenge-notify")

			sentry.CaptureException(errors.New(errorField))
		}
	}()

	todayNotifications := c.getTodayNotifications(challenges)
	if len(todayNotifications) == 0 {
		return
	}

	usersToNotify := make([]notification_case.NewChallengeUser, len(users))
	for i, u := range users {
		usersToNotify[i] = notification_case.NewChallengeUser{
			Email: u.Email,
			Name:  u.Name,
		}
	}

	for i := range todayNotifications {
		c.notificationCase.NewChallenge(context.Internal, notification_case.NewChallengeInput{
			Users:    usersToNotify,
			Schedule: &todayNotifications[i],
		})
	}
}

func (CampaignCase) getTodayNotifications(challenges []CreateScheduledChallenge) []time.Time {
	nowDay, nowMonth, nowYear := time.Now().Date()

	notifications := []time.Time{}
	for _, challenge := range challenges {
		challengeDay, challengeMonth, challengeYear := challenge.StartDate.Date()
		if challengeDay <= nowDay && challengeMonth <= nowMonth && challengeYear <= nowYear {
			notifications = append(notifications, challenge.StartDate)
		}
	}
	return notifications
}

func validateScheduleChallenges(elements []CreateScheduledChallenge) error {
	if len(elements) == 0 {
		return ErrCampaignWithoutChallenges
	}

	uniqueChallenges := make(map[string]bool)
	for _, value := range elements {
		if uniqueChallenges[value.ChallengeID] {
			return ErrCampaignWithDuplicatedChallengesIds
		}
		uniqueChallenges[value.ChallengeID] = true
	}

	startOfToday := utils.StartOfDay(time.Now().UTC())
	endOfToday := utils.EndOfDay(startOfToday)

	for index, el := range elements {
		if el.EndDate.Before(endOfToday) || el.EndDate == endOfToday {
			return ErrPastEndDate(index)
		}

		if el.StartDate.Before(startOfToday) {
			return ErrStartDateBeforeToday(index)
		}

		if el.EndDate.Before(el.StartDate) || el.EndDate == el.StartDate {
			return ErrInvalidRange(index)
		}
	}

	return nil
}

func (u *CampaignCase) getCampaignUsers(
	companyID string,
	onlyDepartments []string,
	onlyUsersCreatedAtGTE time.Time,
	onlyUsersCreatedAtLTE time.Time,
) (*[]department_gateway.GetUsersByFilterGatewayOutput, error) {
	users, err := u.departmentCase.GetUsersByFilter(department_case.GetUsersByFilterCaseInput{
		CompanyID:         companyID,
		Departments:       onlyDepartments,
		UserCompanyStatus: user_company_enum.Active,
		UserCreatedAtGTE:  onlyUsersCreatedAtGTE,
		UserCreatedAtLTE:  onlyUsersCreatedAtLTE,
	})

	if err != nil {
		return nil, err
	}

	if len(*users) == 0 {
		return nil, ErrCampaignWithoutParticipants
	}

	return users, nil
}

func (u *CampaignCase) getQuestions(challenges []CreateScheduledChallenge) (*[]challenge_gateway.GetQuestionsOutput, error) {
	questions, err := u.challengeCase.GetQuestions(parseChallengeIds(challenges))
	if err != nil {
		return nil, err
	}

	if len(*questions) == 0 {
		return nil, ErrCampaignWithoutQuestions
	}

	return questions, nil
}

func makeResults(
	campaignId string,
	questions []challenge_gateway.GetQuestionsOutput,
	participants []g.CreateCampaignParticipant,
) (results []g.Result) {
	for _, question := range questions {
		for _, participant := range participants {
			results = append(results, g.Result{
				UserID:       participant.UserID,
				DepartmentID: participant.DepartmentID,
				ChallengeID:  question.ChallengeID,
				QuestionID:   question.QuestionID,
				CampaignID:   campaignId,
				Tries:        0,
				Correct:      false,
				Status:       campaign_enum.Active,
			})
		}
	}
	return
}

func parseChallengeIds(input []CreateScheduledChallenge) (output []string) {
	for _, element := range input {
		output = append(output, element.ChallengeID)
	}
	return
}

func makeScheduledChallenge(campaignId string, elements []CreateScheduledChallenge) (output []g.ScheduledChallenge) {
	for _, el := range elements {
		output = append(output, g.ScheduledChallenge{
			CampaingID:  campaignId,
			ChallengeID: el.ChallengeID,
			StartDate:   el.StartDate,
			EndDate:     el.EndDate,
		})
	}

	return
}

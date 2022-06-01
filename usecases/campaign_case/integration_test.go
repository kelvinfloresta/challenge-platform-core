package campaign_case_test

import (
	"conformity-core/enums/campaign_enum"
	"conformity-core/enums/user_company_enum"
	"conformity-core/fixtures"
	g "conformity-core/gateways/campaign_gateway"
	"conformity-core/gateways/company_gateway"
	"conformity-core/usecases/campaign_case"
	"conformity-core/utils"
	goErr "errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCampaignCase_Create__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.CampaignCase.NewIntegration()
	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title: "Campaign 1",
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: challenge.ID,
				StartDate:   time.Date(2200, time.January, 0, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2200, time.February, 0, 0, 0, 0, 0, time.UTC),
			},
		},
		CompanyID: user.CompanyID,
	}
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})
	id, err := sut.Create(ctx, campaignToCreate)
	require.Nil(t, err)

	found, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Equal(t, found.ID, id)
	assert.Equal(t, found.Title, campaignToCreate.Title)

	assert.Len(t, found.Results, 1)
	result := found.Results[0]
	assert.Equal(t, user.ID, result.UserID)
	assert.Equal(t, user.DepartmentID, result.DepartmentID)
	assert.Equal(t, challenge.ID, result.ChallengeID)
	assert.Equal(t, uint8(0), result.Tries)
	assert.Equal(t, false, result.Correct)
	assert.Equal(t, campaign_enum.Active, result.Status)
}

func TestCampaignCase_Create__Should_fail_if_company_does_not_have_departments(t *testing.T) {
	fixtures.CleanTestDatabase()

	companyId := fixtures.CreateCompany(t)
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.CampaignCase.NewIntegration()
	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title: "Campaign 1",
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: challenge.ID,
				StartDate:   time.Date(2200, time.January, 0, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2200, time.February, 0, 0, 0, 0, 0, time.UTC),
			},
		},
		CompanyID: companyId,
	}

	_, err := sut.Create(fixtures.FakeBackofficeCtx, campaignToCreate)
	assert.ErrorIs(t, err, campaign_case.ErrCampaignWithoutParticipants)
}

func TestCampaignCase_Create__Should_fail_if_does_not_have_users_in_campaign(t *testing.T) {
	fixtures.CleanTestDatabase()

	companyId := fixtures.CreateCompany(t)
	challenge := fixtures.CreateChallenge(t)
	fixtures.CreateDepartment(t, companyId)

	sut := fixtures.CampaignCase.NewIntegration()
	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title: "Campaign 1",
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: challenge.ID,
				StartDate:   time.Date(2200, time.January, 0, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2200, time.February, 0, 0, 0, 0, 0, time.UTC),
			},
		},
		CompanyID: companyId,
	}

	_, err := sut.Create(fixtures.FakeBackofficeCtx, campaignToCreate)
	assert.Equal(t, err, campaign_case.ErrCampaignWithoutParticipants)
}

func TestCampaignCase_Create__Should_not_add_users_from_departments_of_different_company(t *testing.T) {
	fixtures.CleanTestDatabase()

	COMPANY_WITH_USERS := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, COMPANY_WITH_USERS)
	fixtures.CreateUser(t, &departmentId)

	COMPANY_WITHOUT_USERS := fixtures.CreateCompany(t)
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.CampaignCase.NewIntegration()
	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title: "Campaign 1",
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: challenge.ID,
				StartDate:   time.Date(2200, time.January, 0, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2200, time.February, 0, 0, 0, 0, 0, time.UTC),
			},
		},
		CompanyID: COMPANY_WITHOUT_USERS,
	}

	_, err := sut.Create(fixtures.FakeBackofficeCtx, campaignToCreate)
	assert.ErrorIs(t, err, campaign_case.ErrCampaignWithoutParticipants)
}

func TestCampaignCase_Create__Should_create_schedule_campaign(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)

	campaign, err := sut.GetById(defaultCampaign.CampaignID)
	require.Nil(t, err)

	scheduledChallenges := defaultCampaign.ScheduledChallenge
	for i, challenge := range campaign.ScheduledChallenges {
		require.Equal(t, campaign.ID, challenge.CampaingID)
		require.Equal(t, scheduledChallenges[i].ChallengeID, challenge.ChallengeID)
		require.Equal(t, scheduledChallenges[i].StartDate.Unix(), challenge.StartDate.Unix())
		require.Equal(t, scheduledChallenges[i].EndDate.Unix(), challenge.EndDate.Unix())
	}
}

func TestCampaignCase_Create__Should_fail_if_EndDate_is_before_StartDate(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.CampaignCase.NewIntegration()
	challengesToSchedule := []campaign_case.CreateScheduledChallenge{
		{
			ChallengeID: challenge.ID,
			EndDate:     time.Date(2200, time.February, 24, 0, 0, 0, 0, time.UTC),
			StartDate:   time.Date(2200, time.February, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title:           "Campaign 1",
		Challenges:      challengesToSchedule,
		OnlyDepartments: nil,
		CompanyID:       user.CompanyID,
	}
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})
	_, err := sut.Create(ctx, campaignToCreate)
	require.EqualError(t, err, campaign_case.ErrInvalidRange(0).Error())
}

func TestCampaignCase_Create__Should_fail_if_EndDate_is_equal_to_StartDate(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.CampaignCase.NewIntegration()
	challengesToSchedule := []campaign_case.CreateScheduledChallenge{
		{
			ChallengeID: challenge.ID,
			EndDate:     time.Date(2200, time.February, 24, 0, 0, 0, 0, time.UTC),
			StartDate:   time.Date(2200, time.February, 24, 0, 0, 0, 0, time.UTC),
		},
	}

	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title:           "Campaign 1",
		Challenges:      challengesToSchedule,
		OnlyDepartments: nil,
		CompanyID:       user.CompanyID,
	}
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})
	_, err := sut.Create(ctx, campaignToCreate)
	require.EqualError(t, err, campaign_case.ErrInvalidRange(0).Error())
}

func TestCampaignCase_Create__Should_fail_if_EndDate_is_before_now(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)

	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.CampaignCase.NewIntegration()
	startOfToday := utils.StartOfDay(time.Now().UTC())
	challengesToSchedule := []campaign_case.CreateScheduledChallenge{
		{
			ChallengeID: challenge.ID,
			StartDate:   startOfToday,
			EndDate:     time.Now().UTC().Add(-1),
		},
	}

	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title:           "Campaign 1",
		Challenges:      challengesToSchedule,
		OnlyDepartments: nil,
		CompanyID:       user.CompanyID,
	}
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})
	_, err := sut.Create(ctx, campaignToCreate)
	require.EqualError(t, err, campaign_case.ErrPastEndDate(0).Error())
}

func TestCampaignCase_Create__Should_fail_if_EndDate_is_today(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.CampaignCase.NewIntegration()
	startOfToday := utils.StartOfDay(time.Now().UTC())
	endOfToday := utils.EndOfDay(startOfToday)
	challengesToSchedule := []campaign_case.CreateScheduledChallenge{
		{
			ChallengeID: challenge.ID,
			StartDate:   startOfToday,
			EndDate:     endOfToday,
		},
	}

	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title:           "Campaign 1",
		Challenges:      challengesToSchedule,
		OnlyDepartments: nil,
		CompanyID:       user.CompanyID,
	}
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})
	_, err := sut.Create(ctx, campaignToCreate)
	require.EqualError(t, err, campaign_case.ErrPastEndDate(0).Error())
}
func TestCampaignCase_Create__Should_not_fail_if_StartDate_is_today(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	challenge := fixtures.CreateChallenge(t)

	startOfToday := utils.StartOfDay(time.Now().UTC())
	endDate := startOfToday.AddDate(0, 0, 7)
	sut := fixtures.CampaignCase.NewIntegration()
	challengesToSchedule := []campaign_case.CreateScheduledChallenge{
		{
			ChallengeID: challenge.ID,
			StartDate:   startOfToday,
			EndDate:     endDate,
		},
	}

	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title:           "Campaign 1",
		Challenges:      challengesToSchedule,
		OnlyDepartments: nil,
		CompanyID:       user.CompanyID,
	}
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})
	companyId, err := sut.Create(ctx, campaignToCreate)
	require.Nil(t, err)
	require.NotEmpty(t, companyId)
}

func TestCampaignCase_Create__Should_fail_if_StartDate_is_yesterday(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	challenge := fixtures.CreateChallenge(t)

	yesterday := utils.StartOfDay(time.Now().UTC()).Add(-1)
	endDate := time.Now().UTC().AddDate(0, 0, 7)
	sut := fixtures.CampaignCase.NewIntegration()
	challengesToSchedule := []campaign_case.CreateScheduledChallenge{
		{
			ChallengeID: challenge.ID,
			StartDate:   yesterday,
			EndDate:     endDate,
		},
	}

	campaignToCreate := campaign_case.CreateCampaignCaseInput{
		Title:           "Campaign 1",
		Challenges:      challengesToSchedule,
		OnlyDepartments: nil,
		CompanyID:       user.CompanyID,
	}
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})
	_, err := sut.Create(ctx, campaignToCreate)
	require.EqualError(t, err, campaign_case.ErrStartDateBeforeToday(0).Error())
}

func Test_Create__Should_not_include_suspended_users(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	suspendedUser := fixtures.CreateUser(t, &user.DepartmentID)

	sut := fixtures.CampaignCase.NewIntegration()
	suspended, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		UserID:       suspendedUser.ID,
		DepartmentID: suspendedUser.DepartmentID,
		Status:       user_company_enum.Suspended,
	})

	require.Nil(t, err)
	require.True(t, suspended)

	challenge := fixtures.CreateChallenge(t)
	campaignId, err := sut.Create(ctx, campaign_case.CreateCampaignCaseInput{
		Title: "Campaign 1",
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: challenge.ID,
				StartDate:   time.Now(),
				EndDate:     time.Now().AddDate(0, 0, 1),
			},
		},
		CompanyID: user.CompanyID,
	})
	require.Nil(t, err)
	campaign, err := sut.GetById(campaignId)

	require.Nil(t, err)
	require.Len(t, campaign.Results, 1)
	require.Equal(t, campaign.Results[0].UserID, user.ID)
}

func Test_Create__Should_correctly_filter_by_user_creation_date(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()

	manager := fixtures.CreateCompanyManager(t, nil)

	now := time.Now()
	_, err := sut.Create(manager.Ctx, campaign_case.CreateCampaignCaseInput{
		Title: "Test",
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: fixtures.CreateChallenge(t).ID,
				StartDate:   now,
				EndDate:     now.AddDate(0, 0, 7),
			},
		},
		OnlyUsersCreatedAtGTE: now.Add(-time.Minute),
	})

	require.Nil(t, err)
	users, err := sut.ListUsers(manager.Ctx, campaign_case.ListUsersInput{
		StartDateLTE: time.Time{},
		StartDate:    time.Time{},
		EndDateLTE:   time.Time{},
		EndDateGT:    time.Time{},
	})
	require.Nil(t, err)
	require.Len(t, users, 1)
}

func Test_Create__Should_not_include_users_whose_creation_date_does_not_match(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()

	manager := fixtures.CreateCompanyManager(t, nil)

	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	_, err := sut.Create(manager.Ctx, campaign_case.CreateCampaignCaseInput{
		Title: "Test",
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: fixtures.CreateChallenge(t).ID,
				StartDate:   now,
				EndDate:     now.AddDate(0, 0, 7),
			},
		},
		OnlyUsersCreatedAtGTE: tomorrow,
	})

	require.ErrorIs(t, err, campaign_case.ErrCampaignWithoutParticipants)
}

func Test_AnswerChallenge__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	sut := fixtures.CampaignCase.NewIntegration()

	correctAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	}

	result, err := sut.AnswerChallenge(ctx, correctAnswer)

	require.Nil(t, err)
	require.True(t, result.Correct)
	require.Equal(t, uint8(2), result.RemainingTries)
}

func Test_AnswerChallenge_Should_return_false_if_option_is_wrong(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()
	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	wrongOption := fixtures.GetWrongOption(question)

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	result, err := sut.AnswerChallenge(ctx, wrongAnswer)
	assert.Nil(t, err)
	assert.False(t, result.Correct)
	assert.Equal(t, uint8(2), result.RemainingTries)
}

func Test_AnswerChallenge_Should_fail_if_no_tries_remaining(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()
	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	wrongOption := fixtures.GetWrongOption(question)

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	_, err := sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	_, err = sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	_, err = sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)

	_, err = sut.AnswerChallenge(ctx, wrongAnswer)
	require.ErrorIs(t, err, campaign_case.ErrNoTriesRemaining)
}

func Test_AnswerChallenge_Should_fail_if_no_tries_remaining_even_if_the_answer_is_right(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()
	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	wrongOption := fixtures.GetWrongOption(question)

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	_, err := sut.AnswerChallenge(fixtures.DUMMY_CONTEXT, wrongAnswer)
	require.Nil(t, err)
	_, err = sut.AnswerChallenge(fixtures.DUMMY_CONTEXT, wrongAnswer)
	require.Nil(t, err)
	_, err = sut.AnswerChallenge(fixtures.DUMMY_CONTEXT, wrongAnswer)
	require.Nil(t, err)

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	correctOption := fixtures.GetCorrectOption(question)
	correctAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	}

	_, err = sut.AnswerChallenge(ctx, correctAnswer)
	require.ErrorIs(t, err, campaign_case.ErrNoTriesRemaining)
}

func Test_AnswerChallenge_Should_decrement_tries(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()
	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	wrongOption := fixtures.GetWrongOption(question)

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	result, err := sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	require.Equal(t, uint8(2), result.RemainingTries)
	result, err = sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	require.Equal(t, uint8(1), result.RemainingTries)
	result, err = sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	require.Equal(t, uint8(0), result.RemainingTries)
}

func Test_AnswerChallenge_Should_fail_if_try_to_answer_a_question_that_is_already_correct(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()
	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)
	wrongOption := fixtures.GetWrongOption(question)

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	correctAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	}

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	result, err := sut.AnswerChallenge(ctx, correctAnswer)
	require.Nil(t, err)
	require.Equal(t, result.Correct, true)
	require.Equal(t, result.RemainingTries, uint8(2))

	result, err = sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, result)
	require.ErrorIs(t, err, campaign_case.ErrQuestionIsAlreadyCorrect)
}

func Test_AnswerChallenge_Should_not_alter_tries_if_answered_a_question_that_is_already_correct(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()
	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)
	wrongOption := fixtures.GetWrongOption(question)

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	correctAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	}

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	answer1, err := sut.AnswerChallenge(ctx, correctAnswer)
	require.Nil(t, err)
	require.Equal(t, answer1.Correct, true)
	require.Equal(t, answer1.RemainingTries, uint8(2))

	answer2, err := sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, answer2)
	require.ErrorIs(t, err, campaign_case.ErrQuestionIsAlreadyCorrect)

	result, err := sut.GetResult(ctx, g.GetResultInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CampaignID:   defaultCampaign.CampaignID,
		ChallengeID:  defaultCampaign.Challenge.ID,
		QuestionID:   question.ID,
	})

	require.Nil(t, err)
	require.Equal(t, result.Tries, uint8(1))
	require.Equal(t, result.Correct, true)
}

func Test_AnswerChallenge_Should_be_thread_safe(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()
	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	wrongOption := fixtures.GetWrongOption(question)

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	job, wg := fixtures.AsyncAnswerChallenge(t, ctx, sut)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go job(ctx, wrongAnswer)
	}
	wg.Wait()

	result, err := sut.GetResult(fixtures.DUMMY_CONTEXT, g.GetResultInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CampaignID:   defaultCampaign.CampaignID,
		ChallengeID:  defaultCampaign.Challenge.ID,
		QuestionID:   question.ID,
	})

	require.Nil(t, err)
	require.Equal(t, uint8(3), result.Tries)
}

func Test_AnswerChallenge_Should_not_be_able_to_answer_if_challenge_expired(t *testing.T) {
	t.Skip()
}

func Test_GetDepartmentAVGResult_Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)

	sut := fixtures.CampaignCase.NewIntegration()

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	_, err := sut.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})
	require.Nil(t, err)

	result, err := sut.GetDepartmentAVGResult(ctx, campaign_case.GetDepartmentAVGResultInput{
		CampaignID: defaultCampaign.CampaignID,
	})

	require.Nil(t, err)
	require.Len(t, result, 1)
	require.Equal(t, uint8(100), result[0].AVG)
}

func Test_GetDepartmentAVGResult_Should_return_67_if_missed_once(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)
	wrongOption := fixtures.GetWrongOption(question)

	campaignCase := fixtures.CampaignCase.NewIntegration()

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	_, err := campaignCase.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	})
	require.Nil(t, err)

	_, err = campaignCase.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})
	require.Nil(t, err)

	sut := fixtures.CampaignCase.NewIntegration()
	result, err := sut.GetDepartmentAVGResult(ctx, campaign_case.GetDepartmentAVGResultInput{
		CampaignID: defaultCampaign.CampaignID,
	})

	require.Nil(t, err)
	require.Equal(t, uint8(67), result[0].AVG)
}

func Test_GetDepartmentAVGResult_Should_return_34_if_missed_twice(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)
	wrongOption := fixtures.GetWrongOption(question)

	campaignCase := fixtures.CampaignCase.NewIntegration()

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	_, err := campaignCase.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	_, err = campaignCase.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)

	_, err = campaignCase.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})
	require.Nil(t, err)

	sut := fixtures.CampaignCase.NewIntegration()
	result, err := sut.GetDepartmentAVGResult(ctx, campaign_case.GetDepartmentAVGResultInput{
		CampaignID: defaultCampaign.CampaignID,
	})

	require.Nil(t, err)
	require.Equal(t, uint8(34), result[0].AVG)
}

func Test_GetDepartmentAVGResult_Should_return_0_if_missed_three_times(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	question := defaultCampaign.Challenge.Questions[0]
	wrongOption := fixtures.GetWrongOption(question)

	campaignCase := fixtures.CampaignCase.NewIntegration()

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	_, err := campaignCase.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	_, err = campaignCase.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	_, err = campaignCase.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)

	sut := fixtures.CampaignCase.NewIntegration()
	result, err := sut.GetDepartmentAVGResult(ctx, campaign_case.GetDepartmentAVGResultInput{
		CampaignID: defaultCampaign.CampaignID,
	})

	require.Nil(t, err)
	require.Equal(t, uint8(0), result[0].AVG)
}

func Test_GetDepartmentAVGResult_Should_not_calculate_results_from_non_started_challenges(t *testing.T) {
	fixtures.CleanTestDatabase()
	challenge := fixtures.CreateChallenge(t)
	secondChallenge := fixtures.CreateChallenge(t)
	now := time.Now().UTC()
	tomorrow := now.Add(24 * time.Hour)
	campaign := fixtures.CreateCustomCampaign(t, fixtures.CreateCustomCampaignInput{
		UserAmount: 1,
		Challenges: []campaign_case.CreateScheduledChallenge{
			{ChallengeID: challenge.ID, StartDate: now, EndDate: tomorrow},
			{ChallengeID: secondChallenge.ID, StartDate: tomorrow, EndDate: now.Add(48 * time.Hour)},
		},
	})
	userId := campaign.UserIDs[0]
	question := challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       userId,
		DepartmentID: campaign.DepartmentA,
		CompanyID:    campaign.CompanyID,
	})

	_, err := fixtures.CampaignCase.NewIntegration().AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  campaign.CampaignID,
		ChallengeID: challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})
	require.Nil(t, err)

	sut := fixtures.CampaignCase.NewIntegration()
	result, err := sut.GetDepartmentAVGResult(ctx, campaign_case.GetDepartmentAVGResultInput{
		CampaignID: campaign.CampaignID,
	})

	require.Nil(t, err)
	require.Equal(t, uint8(100), result[0].AVG)
}

func Test_GetDepartmentAVGResult_Should_not_calculate_suspended_results(t *testing.T) {
	fixtures.CleanTestDatabase()
	challengeA := fixtures.CreateChallenge(t)
	challengeB := fixtures.CreateChallenge(t)

	now := time.Now().UTC()
	tomorrow := now.Add(24 * time.Hour)
	campaign := fixtures.CreateCustomCampaign(t, fixtures.CreateCustomCampaignInput{
		UserAmount: 1,
		Challenges: []campaign_case.CreateScheduledChallenge{
			{ChallengeID: challengeA.ID, StartDate: now, EndDate: tomorrow},
			{ChallengeID: challengeB.ID, StartDate: now, EndDate: tomorrow},
		},
	})

	userId := campaign.UserIDs[0]
	question := challengeA.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       userId,
		DepartmentID: campaign.DepartmentA,
	})

	_, err := fixtures.CampaignCase.NewIntegration().AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  campaign.CampaignID,
		ChallengeID: challengeA.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})
	require.Nil(t, err)

	sut := fixtures.CampaignCase.NewIntegration()
	suspended, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		UserID:       userId,
		DepartmentID: campaign.DepartmentA,
		Status:       user_company_enum.Suspended,
	})

	require.True(t, suspended)
	require.Nil(t, err)

	result, err := sut.GetDepartmentAVGResult(ctx, campaign_case.GetDepartmentAVGResultInput{
		CampaignID: campaign.CampaignID,
	})

	require.Nil(t, err)
	require.Len(t, result, 1)
	require.Equal(t, uint8(100), result[0].AVG)
}

func Test_UpdateResult_Should_not_suspended_challenges_already_correct(t *testing.T) {
	fixtures.CleanTestDatabase()
	campaign := fixtures.CreateDefaultCampaign(t, nil)
	question := campaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       campaign.UserID,
		DepartmentID: campaign.DepartmentID,
		CompanyID:    campaign.CompanyID,
	})

	sut := fixtures.CampaignCase.NewIntegration()
	_, err := sut.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  campaign.CampaignID,
		ChallengeID: campaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})
	require.Nil(t, err)

	_, err = sut.UpdateResult(campaign_case.UpdateResultInput{
		UserID:       campaign.UserID,
		DepartmentID: campaign.DepartmentID,
	}, nil)
	require.Nil(t, err)

	result, err := sut.GetResult(ctx, g.GetResultInput{
		UserID:       campaign.UserID,
		DepartmentID: campaign.DepartmentID,
		CampaignID:   campaign.CampaignID,
		ChallengeID:  campaign.Challenge.ID,
		QuestionID:   question.ID,
	})

	require.Nil(t, err)
	require.Equal(t, campaign_enum.Active, result.Status)
}

func Test_UpdateResult_Should_not_suspended_challenges_if_all_attemps_are_over(t *testing.T) {
	fixtures.CleanTestDatabase()
	campaign := fixtures.CreateDefaultCampaign(t, nil)
	question := campaign.Challenge.Questions[0]
	wrongOption := fixtures.GetWrongOption(question)

	sut := fixtures.CampaignCase.NewIntegration()

	wrongAnswer := campaign_case.AnswerChallengeInput{
		CampaignID:  campaign.CampaignID,
		ChallengeID: campaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	}

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       campaign.UserID,
		DepartmentID: campaign.DepartmentID,
		CompanyID:    campaign.CompanyID,
	})

	_, err := sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	_, err = sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)
	_, err = sut.AnswerChallenge(ctx, wrongAnswer)
	require.Nil(t, err)

	_, err = sut.UpdateResult(campaign_case.UpdateResultInput{
		UserID:       campaign.UserID,
		DepartmentID: campaign.DepartmentID,
	}, nil)
	require.Nil(t, err)

	result, err := sut.GetResult(ctx, g.GetResultInput{
		UserID:       campaign.UserID,
		DepartmentID: campaign.DepartmentID,
		CampaignID:   campaign.CampaignID,
		ChallengeID:  campaign.Challenge.ID,
		QuestionID:   question.ID,
	})

	require.Nil(t, err)
	require.Equal(t, campaign_enum.Active, result.Status)
}

func Test_GetUserAVGResult_Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)

	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)
	sut := fixtures.CampaignCase.NewIntegration()

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
	})

	_, err := sut.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})
	require.Nil(t, err)

	result, err := sut.GetUserAVGResult(ctx, campaign_case.GetUserAVGResultInput{
		CampaignID: defaultCampaign.CampaignID,
	})

	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, uint8(100), result.AVG)
	require.Equal(t, uint8(100), result.CompanyAVG)
}

func Test_GetUserAVGResult_Should_not_fail_if_campaign_does_not_exist(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CampaignCase.NewIntegration()
	NON_EXISTENT_CAMPAIGN := fixtures.UUID(t)

	result, err := sut.GetUserAVGResult(fixtures.DUMMY_CONTEXT, campaign_case.GetUserAVGResultInput{
		CampaignID: NON_EXISTENT_CAMPAIGN,
	})

	require.Nil(t, err)
	require.NotNil(t, result)
	require.Zero(t, result.AVG)
	require.Zero(t, result.CompanyAVG)
}

func Test_GetUserAVGResult_Should_calculate_result_from_other_departments(t *testing.T) {
	fixtures.CleanTestDatabase()
	challenge := fixtures.CreateChallenge(t)

	now := time.Now().UTC()
	tomorrow := now.Add(24 * time.Hour)
	campaign := fixtures.CreateCustomCampaign(t, fixtures.CreateCustomCampaignInput{
		UserAmount: 2,
		Challenges: []campaign_case.CreateScheduledChallenge{
			{ChallengeID: challenge.ID, StartDate: now, EndDate: tomorrow},
		},
	})

	question := challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)
	sut := fixtures.CampaignCase.NewIntegration()

	user1Ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       campaign.UserIDs[0],
		DepartmentID: campaign.DepartmentA,
	})

	_, err := sut.AnswerChallenge(user1Ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  campaign.CampaignID,
		ChallengeID: challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})

	require.Nil(t, err)

	user2Ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       campaign.UserIDs[1],
		DepartmentID: campaign.DepartmentB,
	})

	wrongOption := fixtures.GetWrongOption(question)
	for i := 0; i < 3; i++ {
		_, err = sut.AnswerChallenge(user2Ctx, campaign_case.AnswerChallengeInput{
			CampaignID:  campaign.CampaignID,
			ChallengeID: challenge.ID,
			QuestionID:  question.ID,
			OptionID:    wrongOption.ID,
		})
		require.Nil(t, err)
	}

	result, err := sut.GetUserAVGResult(user1Ctx, campaign_case.GetUserAVGResultInput{
		CampaignID: campaign.CampaignID,
	})

	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, uint8(100), result.AVG)
	require.Equal(t, uint8(50), result.CompanyAVG)
}

func Test_ActiveCampaigns__Should_do_happy_path(t *testing.T) {
	t.Skip()
}

func Test_ListChallenges__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)

	sut := fixtures.CampaignCase.NewIntegration()
	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
	})

	isPending := true
	challenges, err := sut.ListChallenges(ctx, campaign_case.ListChallengesInput{
		IsPending: &isPending,
	})

	require.Nil(t, err)
	require.Len(t, challenges, 1)

	for _, found := range challenges {
		scheduledChallenge := defaultCampaign.ScheduledChallenge[0]
		require.Equal(t, scheduledChallenge.ChallengeID, found.ID)
		require.Equal(t, scheduledChallenge.StartDate.Unix(), found.StartDate.Unix())
		require.Equal(t, scheduledChallenge.EndDate.Unix(), found.EndDate.Unix())
	}
}

func Test_ListChallenges__Should_not_list_completed_challenges_when_listing_pending_ones(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)

	sut := fixtures.CampaignCase.NewIntegration()
	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
	})

	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)

	_, err := sut.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})

	require.Nil(t, err)

	IsPending := true
	challenges, err := sut.ListChallenges(ctx, campaign_case.ListChallengesInput{
		IsPending: &IsPending,
	})

	require.Nil(t, err)
	require.Len(t, challenges, 0)
}

func Test_ListChallenges__Should_list_if_the_user_answered_a_question_but_not_all(t *testing.T) {
	t.Skip()
}

func Test_ListChallenges__Should_never_list_uninitiated_challenges(t *testing.T) {
	challenge := fixtures.CreateChallenge(t)

	now := time.Now()
	customCampaign := fixtures.CreateCustomCampaign(t, fixtures.CreateCustomCampaignInput{
		UserAmount: 1,
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: challenge.ID,
				StartDate:   now.AddDate(0, 0, 1),
				EndDate:     now.AddDate(0, 0, 7),
			},
		},
	})

	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       customCampaign.UserIDs[0],
		DepartmentID: customCampaign.DepartmentA,
	})

	sut := fixtures.CampaignCase.NewIntegration()

	notIsPending := false
	challenges, err := sut.ListChallenges(ctx, campaign_case.ListChallengesInput{
		IsPending: &notIsPending,
	})

	require.Nil(t, err)
	require.Len(t, challenges, 0)

	isPending := true
	challenges, err = sut.ListChallenges(ctx, campaign_case.ListChallengesInput{
		IsPending: &isPending,
	})

	require.Nil(t, err)
	require.Len(t, challenges, 0)
}

func Test_ListChallenges__Should_be_able_to_list_non_pending_challenges(t *testing.T) {
	challenge := fixtures.CreateChallenge(t)

	now := time.Now()
	customCampaign := fixtures.CreateCustomCampaign(t, fixtures.CreateCustomCampaignInput{
		UserAmount: 1,
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: challenge.ID,
				StartDate:   now.AddDate(0, 0, 1),
				EndDate:     now.AddDate(0, 0, 7),
			},
		},
	})

	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       customCampaign.UserIDs[0],
		DepartmentID: customCampaign.DepartmentA,
	})

	sut := fixtures.CampaignCase.NewIntegration()

	isPending := true
	challenges, err := sut.ListChallenges(ctx, campaign_case.ListChallengesInput{
		IsPending: &isPending,
	})

	require.Nil(t, err)
	require.Len(t, challenges, 0)

	challenges, err = sut.ListChallenges(ctx, campaign_case.ListChallengesInput{
		IsPending: &isPending,
	})

	require.Nil(t, err)
	require.Len(t, challenges, 0)
}
func Test_ListChallenges__Should_list_only_challenges_from_own_company(t *testing.T) {
	t.Skip()
}

func Test_ListQuestions__Should_do_happy_path(t *testing.T) {
	t.Skip()
}

func Test_ListQuestions__Should_not_list_if_correct_is_true(t *testing.T) {
	t.Skip()
}

func Test_ListQuestions__Should_not_list_if_tries_is_less_than_maxTries(t *testing.T) {
	t.Skip()
}

func Test_List__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)

	question := defaultCampaign.Challenge.Questions[0]
	correctOption := fixtures.GetCorrectOption(question)
	sut := fixtures.CampaignCase.NewIntegration()

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	_, err := sut.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  defaultCampaign.CampaignID,
		ChallengeID: defaultCampaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})
	require.Nil(t, err)

	result, err := sut.List(ctx)
	require.Nil(t, err)
	require.Equal(t, len(result), 1)

}

func Test_ListUsers__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()
	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	sut := fixtures.CampaignCase.NewIntegration()
	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		CompanyID:    defaultCampaign.CompanyID,
	})

	challenges, err := sut.ListUsers(ctx, campaign_case.ListUsersInput{
		StartDate:  time.Now(),
		EndDateGT:  defaultCampaign.EndDate.Add(-1 * time.Second),
		EndDateLTE: defaultCampaign.EndDate.AddDate(0, 0, 1),
	})

	require.Nil(t, err)

	require.Len(t, challenges, 1)
}

func Test_ChangeUserStatus__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)

	sut := fixtures.CampaignCase.NewIntegration()

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
	})

	changed, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
		Status:       user_company_enum.Suspended,
	})

	require.Nil(t, err)
	require.True(t, changed)

	companyCase := fixtures.CompanyCase.NewIntegration()
	user, err := companyCase.GetUser(company_gateway.GetUserInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
	})

	require.Nil(t, err)
	require.Equal(t, user_company_enum.Suspended, user.Status)
}

func Test_ChangeUserStatus__Should_be_able_to_change_status_from_suspended_to_active(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	sut := fixtures.CampaignCase.NewIntegration()

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
	})

	changed, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
		Status:       user_company_enum.Suspended,
	})

	require.Nil(t, err)
	require.True(t, changed)

	user, err := fixtures.CompanyCase.NewIntegration().GetUser(company_gateway.GetUserInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
	})

	require.Nil(t, err)
	require.Equal(t, user_company_enum.Suspended, user.Status)

	changed, err = sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
		Status:       user_company_enum.Active,
	})

	require.Nil(t, err)
	require.True(t, changed)

	user, err = fixtures.CompanyCase.NewIntegration().GetUser(company_gateway.GetUserInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
	})

	require.Nil(t, err)
	require.Equal(t, user_company_enum.Active, user.Status)
}

func Test_ChangeUserStatus__Should_fail_if_user_does_not_exist(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)

	sut := fixtures.CampaignCase.NewIntegration()
	NON_EXISTENT_USER := fixtures.UUID(t)

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
	})

	changed, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       NON_EXISTENT_USER,
		Status:       user_company_enum.Suspended,
	})

	require.Nil(t, err)
	require.False(t, changed)
}

func Test_ChangeUserStatus__Should_return_false_if_user_is_not_in_campaign(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)

	sut := fixtures.CampaignCase.NewIntegration()
	NON_EXISTENT_USER := fixtures.UUID(t)

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
	})

	changed, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       NON_EXISTENT_USER,
		Status:       user_company_enum.Suspended,
	})

	require.Nil(t, err)
	require.False(t, changed)
}

func Test_ChangeUserStatus__Should_suspend_result_if_new_status_is_Suspended(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	sut := fixtures.CampaignCase.NewIntegration()
	NEW_STATUS := user_company_enum.Suspended

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
	})

	changed, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
		Status:       NEW_STATUS,
	})

	require.Nil(t, err)
	require.True(t, changed)

	campaignCase := fixtures.CampaignCase.NewIntegration()
	result, err := campaignCase.GetResult(ctx, g.GetResultInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		ChallengeID:  defaultCampaign.Challenge.ID,
		CampaignID:   defaultCampaign.CampaignID,
		QuestionID:   defaultCampaign.Challenge.Questions[0].ID,
	})

	require.Nil(t, err)
	require.Equal(t, campaign_enum.Suspended, result.Status)
}

func Test_ChangeUserStatus__Should_active_results_if_new_status_is_active(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	sut := fixtures.CampaignCase.NewIntegration()

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
	})

	changed, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
		Status:       user_company_enum.Suspended,
	})
	require.Nil(t, err)
	require.True(t, changed)

	campaignCase := fixtures.CampaignCase.NewIntegration()
	result, err := campaignCase.GetResult(ctx, g.GetResultInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		ChallengeID:  defaultCampaign.Challenge.ID,
		CampaignID:   defaultCampaign.CampaignID,
		QuestionID:   defaultCampaign.Challenge.Questions[0].ID,
	})
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, campaign_enum.Suspended, result.Status)

	changed, err = sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
		Status:       user_company_enum.Active,
	})
	require.Nil(t, err)
	require.True(t, changed)

	result, err = campaignCase.GetResult(ctx, g.GetResultInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
		ChallengeID:  defaultCampaign.Challenge.ID,
		CampaignID:   defaultCampaign.CampaignID,
		QuestionID:   defaultCampaign.Challenge.Questions[0].ID,
	})

	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, campaign_enum.Active, result.Status)
}

func Test_ChangeUserStatus__Should_rollback_if_UpdateResult_fail(t *testing.T) {
	fixtures.CleanTestDatabase()

	defaultCampaign := fixtures.CreateDefaultCampaign(t, nil)
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
	})

	gatewayMock := &fixtures.CampaignGatewayMock{}
	gateway := &g.GormCampaignGatewayFacade{DB: fixtures.DB_Test}
	gatewayMock.ChangeUserStatus__Impl = gateway.ChangeUserStatus
	gatewayMock.Transaction__Impl = gateway.Transaction
	gatewayMock.UpdateResult__Error = goErr.New("ANY_ERROR")

	notificationCase, _ := fixtures.NotificationCase.NewUnit()

	sut := campaign_case.New(
		gatewayMock,
		fixtures.DepartmentCase.NewIntegration(),
		fixtures.CompanyCase.NewIntegration(),
		fixtures.ChallengeCase.NewIntegration(),
		notificationCase,
	)

	changed, err := sut.ChangeUserStatus(ctx, campaign_case.ChangeUserStatusCaseInput{
		DepartmentID: defaultCampaign.DepartmentID,
		UserID:       defaultCampaign.UserID,
		Status:       user_company_enum.Suspended,
	})

	require.ErrorIs(t, err, gatewayMock.UpdateResult__Error)
	require.False(t, changed)

	user, err := fixtures.CompanyCase.NewIntegration().GetUser(company_gateway.GetUserInput{
		UserID:       defaultCampaign.UserID,
		DepartmentID: defaultCampaign.DepartmentID,
	})

	require.Nil(t, err)
	require.Equal(t, user_company_enum.Active, user.Status)
}

func Test_RemoveUser__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()
	user := fixtures.CreateUser(t, nil)

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	sut := fixtures.CampaignCase.NewIntegration()
	removed, err := sut.RemoveUser(ctx, campaign_case.RemoveUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.True(t, removed)
}

func Test_RemoveUser__Should_return_false_if_not_found(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)

	sut := fixtures.CampaignCase.NewIntegration()

	NON_EXISTENT_USER := fixtures.UUID(t)
	removed, err := sut.RemoveUser(fixtures.FakeBackofficeCtx, campaign_case.RemoveUserInput{
		UserID:       NON_EXISTENT_USER,
		DepartmentID: departmentId,
	})

	require.Nil(t, err)
	require.False(t, removed)
}

func Test_RemoveUser__Should_rollback_if_UpdateResult_fail(t *testing.T) {
	fixtures.CleanTestDatabase()
	user := fixtures.CreateUser(t, nil)

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	gatewayMock := &fixtures.CampaignGatewayMock{}
	gateway := &g.GormCampaignGatewayFacade{DB: fixtures.DB_Test}
	gatewayMock.RemoveUser__Impl = gateway.RemoveUser
	gatewayMock.Transaction__Impl = gateway.Transaction
	gatewayMock.UpdateResult__Error = goErr.New("ANY_ERROR")

	notificationCase, _ := fixtures.NotificationCase.NewUnit()

	sut := campaign_case.New(
		gatewayMock,
		fixtures.DepartmentCase.NewIntegration(),
		fixtures.CompanyCase.NewIntegration(),
		fixtures.ChallengeCase.NewIntegration(),
		notificationCase,
	)

	removed, err := sut.RemoveUser(ctx, campaign_case.RemoveUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.ErrorIs(t, err, gatewayMock.UpdateResult__Error)
	require.False(t, removed)

	foundUser, err := fixtures.CompanyCase.NewIntegration().GetUser(company_gateway.GetUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, user_company_enum.Active, foundUser.Status)
}

func Test_RemoveUser__Should_not_remove_if_the_department_does_not_belong_to_the_company(t *testing.T) {
	fixtures.CleanTestDatabase()
	user := fixtures.CreateUser(t, nil)
	USER_FROM_ANOTHER_COMPANY := fixtures.CreateUser(t, nil)

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       USER_FROM_ANOTHER_COMPANY.ID,
		DepartmentID: USER_FROM_ANOTHER_COMPANY.DepartmentID,
	})

	sut := fixtures.CampaignCase.NewIntegration()
	removed, err := sut.RemoveUser(ctx, campaign_case.RemoveUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.False(t, removed)

	foundUser, err := fixtures.CompanyCase.NewIntegration().GetUser(company_gateway.GetUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, user_company_enum.Active, foundUser.Status)
}

func Test_RemoveUser__Should_remove_if_the_department_does_not_belong_to_the_company_but_the_role_is_backoffice(t *testing.T) {
	fixtures.CleanTestDatabase()
	user := fixtures.CreateUser(t, nil)
	backoffice := fixtures.CreateUser(t, nil)
	backofficeCtx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       backoffice.ID,
		DepartmentID: backoffice.DepartmentID,
		CompanyID:    backoffice.CompanyID,
		Role:         user_company_enum.RoleBackoffice,
	})

	sut := fixtures.CampaignCase.NewIntegration()

	removed, err := sut.RemoveUser(backofficeCtx, campaign_case.RemoveUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.True(t, removed)

	foundUser, err := fixtures.CompanyCase.NewIntegration().GetUser(company_gateway.GetUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.Nil(t, foundUser)
}

func Test_NotifyNewChallenge__Should_schedule_a_message_at_12_AM_UTC(t *testing.T) {
	fixtures.CleanTestDatabase()

	campaign := fixtures.CreateDefaultCampaign(t, nil)

	notificationCase, notification := fixtures.NotificationCase.NewUnit()
	sut := campaign_case.New(
		&g.GormCampaignGatewayFacade{DB: fixtures.DB_Test},
		fixtures.DepartmentCase.NewIntegration(),
		fixtures.CompanyCase.NewIntegration(),
		fixtures.ChallengeCase.NewIntegration(),
		notificationCase,
	)

	ONE_DAY_BEFORE_ENDING := campaign.EndDate.AddDate(0, 0, -1)
	err := sut.NotifyDeadlineChallenge(ONE_DAY_BEFORE_ENDING)
	require.Nil(t, err)
	require.Len(t, notification.Messages, 1)

	expectedSchedule := utils.StartOfDay(campaign.EndDate).Add(12 * time.Hour)
	for _, message := range notification.Messages {
		require.WithinDuration(t, expectedSchedule, *message.Schedule, 0)
	}
}

func Test_NotifyNewChallenge__Should_calculate_the_day_correctly_even_if_the_day_is_different_because_of_the_time_zone(t *testing.T) {
	fixtures.CleanTestDatabase()

	campaign := fixtures.CreateDefaultCampaign(t, nil)

	notificationCase, notification := fixtures.NotificationCase.NewUnit()

	endDate := campaign.EndDate

	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.Nil(t, err)

	sut := campaign_case.New(
		&g.GormCampaignGatewayFacade{DB: fixtures.DB_Test},
		fixtures.DepartmentCase.NewIntegration(),
		fixtures.CompanyCase.NewIntegration(),
		fixtures.ChallengeCase.NewIntegration(),
		notificationCase,
	)

	ONE_DAY_BEFORE_ENDING := utils.StartOfDay(endDate).AddDate(0, 0, -1)
	DIFF_DAY_BECAUSE_TIME_ZONE := ONE_DAY_BEFORE_ENDING.In(loc)
	err = sut.NotifyDeadlineChallenge(DIFF_DAY_BECAUSE_TIME_ZONE)
	require.Nil(t, err)
	require.Len(t, notification.Messages, 1)

	expectedSchedule := utils.StartOfDay(endDate).Add(12 * time.Hour)
	for _, message := range notification.Messages {
		require.WithinDuration(t, expectedSchedule, *message.Schedule, 0)
	}
}

func Test_NotifyNewChallenge__Should_not_schedule_a_message_if_there_is_more_than_24_hours_to_the_end(t *testing.T) {
	fixtures.CleanTestDatabase()

	campaign := fixtures.CreateDefaultCampaign(t, nil)

	notificationCase, notification := fixtures.NotificationCase.NewUnit()
	sut := campaign_case.New(
		&g.GormCampaignGatewayFacade{DB: fixtures.DB_Test},
		fixtures.DepartmentCase.NewIntegration(),
		fixtures.CompanyCase.NewIntegration(),
		fixtures.ChallengeCase.NewIntegration(),
		notificationCase,
	)

	endDate := campaign.EndDate
	ONE_SECOND_LEFT_TO_SCHEDULE := utils.StartOfDay(endDate.AddDate(0, 0, -1)).Add(-1 * time.Second)
	err := sut.NotifyDeadlineChallenge(ONE_SECOND_LEFT_TO_SCHEDULE)
	require.Nil(t, err)
	require.Len(t, notification.Messages, 0)
}

func Test_NotifyNewChallenge__Should_not_send_a_message_if_there_are_no_users(t *testing.T) {
	fixtures.CleanTestDatabase()

	notificationCase, notification := fixtures.NotificationCase.NewUnit()
	sut := campaign_case.New(
		&g.GormCampaignGatewayFacade{DB: fixtures.DB_Test},
		fixtures.DepartmentCase.NewIntegration(),
		fixtures.CompanyCase.NewIntegration(),
		fixtures.ChallengeCase.NewIntegration(),
		notificationCase,
	)

	now := time.Now()
	err := sut.NotifyDeadlineChallenge(now)
	require.Nil(t, err)
	require.Len(t, notification.Messages, 0)
}

func Test_NotifyNewChallenge__Should_anticipate_notification_if_it_end_before_12_AM(t *testing.T) {
	fixtures.CleanTestDatabase()

	startDate := time.Date(2300, time.March, 3, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2300, time.March, 10, 0, 0, 0, 0, time.UTC)

	fixtures.CreateCustomCampaign(t, fixtures.CreateCustomCampaignInput{
		UserAmount: 1,
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: fixtures.CreateChallenge(t).ID,
				StartDate:   startDate,
				EndDate:     endDate,
			},
		},
	})

	notificationCase, notification := fixtures.NotificationCase.NewUnit()

	sut := campaign_case.New(
		&g.GormCampaignGatewayFacade{DB: fixtures.DB_Test},
		fixtures.DepartmentCase.NewIntegration(),
		fixtures.CompanyCase.NewIntegration(),
		fixtures.ChallengeCase.NewIntegration(),
		notificationCase,
	)

	schedulerTime := time.Date(
		endDate.Year(),
		endDate.Month(),
		endDate.Day()-2,
		23, 0, 0, 0, time.UTC)

	err := sut.NotifyDeadlineChallenge(schedulerTime)
	require.Nil(t, err)
	require.Len(t, notification.Messages, 1)

	expectedSchedule := utils.StartOfDay(endDate.AddDate(0, 0, -1)).Add(12 * time.Hour)
	for _, message := range notification.Messages {
		require.WithinDuration(t, expectedSchedule, *message.Schedule, 0)
	}
}

func Test_ListResults__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	campaign := fixtures.CreateDefaultCampaign(t, nil)

	sut := fixtures.CampaignCase.NewIntegration()

	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       campaign.UserID,
		DepartmentID: campaign.DepartmentID,
	})

	results, err := sut.ListResults(ctx, campaign_case.ListResultsInput{})
	require.Nil(t, err)
	require.Len(t, results, 1)

	for _, r := range results {
		require.Equal(t, campaign.UserID, r.UserID)
		require.Equal(t, campaign.Challenge.ID, r.ChallengeID)
		require.Equal(t, campaign.DepartmentID, r.DepartmentID)
		require.NotEmpty(t, r.Name)
		require.NotEmpty(t, r.DepartmentName)
		require.Zero(t, r.Tries)
		require.False(t, r.Correct)
		require.Equal(t, campaign.Challenge.Title, r.ChallengeTitle)
		require.WithinDuration(t, campaign.EndDate, r.EndDate, time.Second)
		require.False(t, r.Finished)
		require.Zero(t, r.IMPD)
	}
}

func Test_ListResults__Should_correctly_list_answered_challenges(t *testing.T) {
	fixtures.CleanTestDatabase()

	campaign := fixtures.CreateDefaultCampaign(t, nil)

	sut := fixtures.CampaignCase.NewIntegration()
	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       campaign.UserID,
		DepartmentID: campaign.DepartmentID,
	})
	question := campaign.Challenge.Questions[0]
	_, err := sut.AnswerChallenge(ctx, campaign_case.AnswerChallengeInput{
		CampaignID:  campaign.CampaignID,
		ChallengeID: campaign.Challenge.ID,
		QuestionID:  question.ID,
		OptionID:    fixtures.GetCorrectOption(question).ID,
	})
	require.Nil(t, err)

	results, err := sut.ListResults(ctx, campaign_case.ListResultsInput{})
	require.Nil(t, err)
	require.Len(t, results, 1)

	for _, r := range results {
		require.Equal(t, campaign.UserID, r.UserID)
		require.Equal(t, campaign.Challenge.ID, r.ChallengeID)
		require.Equal(t, campaign.DepartmentID, r.DepartmentID)
		require.NotEmpty(t, r.Name)
		require.NotEmpty(t, r.DepartmentName)
		require.Equal(t, uint8(1), r.Tries)
		require.True(t, r.Correct)
		require.Equal(t, campaign.Challenge.Title, r.ChallengeTitle)
		require.WithinDuration(t, campaign.EndDate, r.EndDate, time.Second)
		require.True(t, r.Finished)
		require.Equal(t, uint8(100), r.IMPD)
	}
}

func Test_ListResults__Should_return_an_empty_slice_if_it_has_no_results(t *testing.T) {
	user := fixtures.CreateUser(t, nil)
	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	sut := fixtures.CampaignCase.NewIntegration()
	results, err := sut.ListResults(ctx, campaign_case.ListResultsInput{})
	require.Nil(t, err)
	require.NotNil(t, results)
	require.Empty(t, results)
}

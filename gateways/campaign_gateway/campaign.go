package campaign_gateway

import (
	"conformity-core/context"
	"conformity-core/enums/campaign_enum"
	"conformity-core/enums/user_company_enum"
	database "conformity-core/frameworks/database/gorm"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CreateCampaignParticipant struct {
	UserID       string
	DepartmentID string
}

type Result struct {
	UserID       string
	DepartmentID string
	ChallengeID  string
	CampaignID   string
	QuestionID   string
	Status       campaign_enum.ResultStatus
	Tries        uint8
	Correct      bool
}

type CreateCampaignGatewayInput struct {
	CampaignID          string
	Title               string
	Results             []Result
	ScheduledChallenges []ScheduledChallenge
}

type GetByIdCampaignGatewayOutput struct {
	ID                  string
	Title               string
	Results             []Result
	ScheduledChallenges []*ScheduledChallenge

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type UpdateResultInput struct {
	NewStatus  campaign_enum.ResultStatus
	NewTries   *uint8
	NewCorrect *bool

	EndDateGreaterThan  *time.Time
	ActualTriesLessThan *uint8
	ChallengeIDs        []string
	ActualCorrect       *bool
	ActualTries         *uint8
	CampaignID          string
	ChallengeID         string
	QuestionID          string
	UserID              string
	DepartmentID        string
}

type GetResultInput struct {
	UserID       string
	DepartmentID string
	ChallengeID  string
	CampaignID   string
	QuestionID   string
}

type ScheduledChallenge struct {
	CampaingID  string
	ChallengeID string
	StartDate   time.Time
	EndDate     time.Time
}

type ResultGroupBy string

const (
	GroupByDepartment = "department_id"
	GroupByUser       = "user_id"
)

type GetAVGResultOR struct {
	Correct    *bool
	Tries      *uint8
	EndDateLTE time.Time
}

type GetAVGResultInput struct {
	UserID       string
	DepartmentID string
	CampaignID   string
	CompanyID    string
	Status       campaign_enum.ResultStatus
	GroupBy      ResultGroupBy
	OR           *GetAVGResultOR
}

type GetAVGResultOutput struct {
	AVG            uint8
	DepartmentID   string
	DepartmentName string
}

type GetActiveCampaignsOutput struct {
	Title      string
	CampaignID string
}

type GetChallengesCampaignOutput struct {
	ID         string
	Title      string
	CampaignID string
	IMPD       uint8
	Media      datatypes.JSON
	Questions  datatypes.JSON

	StartDate time.Time
	EndDate   time.Time
}

type ListQuestionsInput struct {
	ChallengeID   string
	CampaignID    string
	UserID        string
	DepartmentID  string
	Correct       *bool
	Status        campaign_enum.ResultStatus
	TriesLessThan *uint8
}

type Option struct {
	ID      string
	Title   string
	Correct bool
}

type ListQuestionsOutput struct {
	ID      string
	Title   string
	Options []Option
	Tries   uint8
}

type ListResultsOutput struct {
	UserID         string                   `json:"user_id"`
	ChallengeID    string                   `json:"challenge_id"`
	DepartmentID   string                   `json:"department_id"`
	Name           string                   `json:"userName"`
	DepartmentName string                   `json:"departmentName"`
	Tries          uint8                    `json:"tries"`
	Correct        bool                     `json:"correct"`
	ChallengeTitle string                   `json:"challenge_title"`
	EndDate        time.Time                `json:"end_date"`
	Status         user_company_enum.Status `json:"string"`
}

type ListOutput struct {
	Title      string
	CampaignID string
	StartDate  time.Time
	EndDate    time.Time
}

type ListInput struct {
	CompanyID string
}

type ListUsersInput struct {
	StartDateLTE         time.Time
	EndDateLTE           time.Time
	EndDateGT            time.Time
	TriesLessThan        *uint8
	Correct              *bool
	UseCompanyStatus     *user_company_enum.Status
	CampaignResultStatus *campaign_enum.ResultStatus
	StartDate            time.Time
	ChallengeID          string
	CampaignID           string
}

type ListUsersOutput struct {
	Email          string
	Name           string
	DepartmentName string
	UserID         string
	DepartmentID   string
	CompanyID      string
	IMPD           uint8
}

type RemoveUserInput struct {
	DepartmentID string
	UserID       string
}

type ChangeUserStatusInput struct {
	DepartmentID string
	UserID       string
	Status       user_company_enum.Status
}

type ListChallengesORInput struct {
	Tries      *uint8
	Correct    *bool
	EndDateLTE time.Time
}

type ListChallengesInput struct {
	EndDateGTE   time.Time
	StartDateLTE time.Time
	UserID       string
	CompanyID    string
	Status       campaign_enum.ResultStatus
	TriesLT      *uint8
	Correct      *bool
	OR           *ListChallengesORInput
}

type ListResultsInput struct {
	UserID      string
	CampaignID  string
	ChallengeID string
	CompanyID   string
	Status      campaign_enum.ResultStatus
}

type CreateResultsInput struct {
	UserID       string
	DepartmentID string
	ChallengeID  string
	CampaignID   string
	QuestionID   string
	Status       campaign_enum.ResultStatus
	Tries        uint8
	Correct      bool
}

type ICampaignGateway interface {
	database.Transaction
	Create(input CreateCampaignGatewayInput) (string, error)
	CreateResults(input []CreateResultsInput, tx *gorm.DB) error
	GetById(id string) (*GetByIdCampaignGatewayOutput, error)
	UpdateResult(input UpdateResultInput, tx *gorm.DB) (bool, error)
	GetResult(ctx *context.CoreCtx, input GetResultInput) (*Result, error)
	GetAVGResult(ctx *context.CoreCtx, input GetAVGResultInput) ([]*GetAVGResultOutput, error)
	ListChallenges(ctx *context.CoreCtx, input ListChallengesInput) ([]GetChallengesCampaignOutput, error)
	ListQuestions(ctx *context.CoreCtx, input ListQuestionsInput) ([]ListQuestionsOutput, error)
	ListResults(ctx *context.CoreCtx, input ListResultsInput) ([]ListResultsOutput, error)
	List(ctx *context.CoreCtx, input ListInput) ([]*ListOutput, error)
	ListUsers(ctx *context.CoreCtx, input ListUsersInput) ([]ListUsersOutput, error)
	RemoveUser(ctx *context.CoreCtx, input RemoveUserInput, tx *gorm.DB) (bool, error)
	ChangeUserStatus(input ChangeUserStatusInput, tx *gorm.DB) (bool, error)
}

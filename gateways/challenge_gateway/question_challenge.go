package challenge_gateway

type GetQuestionsOutput struct {
	ChallengeID string
	QuestionID  string
}

type CreateQuestion struct {
	Title   string
	Options []CreateOption
}

type Question struct {
	ID      string
	Title   string
	Options []Option
}

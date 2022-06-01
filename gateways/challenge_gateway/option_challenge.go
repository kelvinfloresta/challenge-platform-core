package challenge_gateway

type GetOptionInput struct {
	ChallengeID string
	QuestionID  string
	OptionID    string
}

type CreateOption struct {
	Title   string
	Correct bool
}

type Option struct {
	ID      string
	Title   string
	Correct bool
}

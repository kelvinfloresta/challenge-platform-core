package challenge_gateway

type CreateMedia struct {
	Title       string
	Description string
	Path        string
}

type Media struct {
	ID          string
	Title       string
	Description string
	Path        string
}

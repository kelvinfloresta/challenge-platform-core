package domain

type result struct {
	tries   uint8
	correct bool
}

func Result(tries uint8, correct bool) result {
	return result{tries, correct}
}

func (r result) IsFinished() bool {
	return r.correct || r.tries == 3
}

func (r result) IMPD() uint8 {
	if !r.correct {
		return 0
	}

	switch r.tries {
	case 1:
		return 100
	case 2:
		return 67
	case 3:
		return 34
	default:
		return 0
	}
}

package campaign_gateway

import "conformity-core/frameworks/database/gorm/models"

func adaptCreateResult(results []Result) (adaptedResults []*models.CampaignResult) {
	for _, result := range results {
		adaptedResults = append(adaptedResults, &models.CampaignResult{
			CampaignID:   result.CampaignID,
			UserID:       result.UserID,
			DepartmentID: result.DepartmentID,
			ChallengeID:  result.ChallengeID,
			QuestionID:   result.QuestionID,
			Tries:        result.Tries,
			Correct:      result.Correct,
			Status:       result.Status,
		})
	}

	return
}

func adaptResult(results []models.CampaignResult) (adapted []Result) {
	for _, result := range results {
		adapted = append(adapted, Result{
			UserID:       result.UserID,
			DepartmentID: result.DepartmentID,
			ChallengeID:  result.ChallengeID,
			Tries:        result.Tries,
			Correct:      result.Correct,
			Status:       result.Status,
			QuestionID:   result.QuestionID,
			CampaignID:   result.CampaignID,
		})
	}

	return
}

func adaptCreateScheduledChallenge(elements []ScheduledChallenge) (adaptedResults []*models.ScheduledChallenge) {
	for _, scheduledChallenge := range elements {
		adaptedResults = append(adaptedResults, &models.ScheduledChallenge{
			CampaignID:  scheduledChallenge.CampaingID,
			ChallengeID: scheduledChallenge.ChallengeID,
			StartDate:   scheduledChallenge.StartDate,
			EndDate:     scheduledChallenge.EndDate,
		})
	}

	return
}

func adaptScheduledChallenges(model []models.ScheduledChallenge) (adapted []*ScheduledChallenge) {
	for _, sc := range model {
		adapted = append(adapted, &ScheduledChallenge{
			CampaingID:  sc.CampaignID,
			ChallengeID: sc.ChallengeID,
			StartDate:   sc.StartDate,
			EndDate:     sc.EndDate,
		})
	}
	return
}

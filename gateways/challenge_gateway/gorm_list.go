package challenge_gateway

import "conformity-core/frameworks/database/gorm/models"

func (g GormChallengeGatewayFacade) List() (out []ListOutput, err error) {
	result := g.DB.Conn.Model(&models.Challenge{}).Select(`
		challenges.id,
		challenges.title,
		challenges.segment,
		challenges.media->>'Path' as media,
		challenges.media->>'Description' as description
	`).Scan(&out)

	return out, result.Error
}

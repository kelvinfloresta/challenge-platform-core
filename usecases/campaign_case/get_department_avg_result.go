package campaign_case

import (
	"conformity-core/context"
	"conformity-core/enums/campaign_enum"
	g "conformity-core/gateways/campaign_gateway"
	"time"
)

type GetDepartmentAVGResultInput struct {
	CampaignID string
}

func (c CampaignCase) GetDepartmentAVGResult(ctx *context.CoreCtx, input GetDepartmentAVGResultInput) ([]*g.GetAVGResultOutput, error) {
	now := time.Now()
	correct := true
	tries := uint8(3)
	OR := &g.GetAVGResultOR{
		Correct:    &correct,
		Tries:      &tries,
		EndDateLTE: now,
	}

	return c.gateway.GetAVGResult(ctx, g.GetAVGResultInput{
		CampaignID: input.CampaignID,
		CompanyID:  ctx.Session.CompanyID,
		GroupBy:    g.GroupByDepartment,
		Status:     campaign_enum.Active,
		OR:         OR,
	})
}

package campaign_case

import (
	"conformity-core/context"
	"conformity-core/enums/campaign_enum"
	g "conformity-core/gateways/campaign_gateway"
	"time"
)

type GetUserAVGResultInput struct {
	CampaignID string
}

type GetUserAVGResultOutput struct {
	AVG        uint8 `json:"avg"`
	CompanyAVG uint8 `json:"companyAvg"`
}

func (c CampaignCase) GetUserAVGResult(ctx *context.CoreCtx, input GetUserAVGResultInput) (out GetUserAVGResultOutput, err error) {
	now := time.Now()
	correct := true
	tries := uint8(3)
	OR := &g.GetAVGResultOR{
		Correct:    &correct,
		Tries:      &tries,
		EndDateLTE: now,
	}

	userAvg, err := c.gateway.GetAVGResult(ctx, g.GetAVGResultInput{
		CampaignID:   input.CampaignID,
		UserID:       ctx.Session.UserID,
		DepartmentID: ctx.Session.DepartmentID,
		Status:       campaign_enum.Active,
		OR:           OR,
	})

	if err != nil {
		return
	}

	if len(userAvg) != 0 {
		out.AVG = userAvg[0].AVG
	}

	companyAvg, err := c.gateway.GetAVGResult(ctx, g.GetAVGResultInput{
		CompanyID:  ctx.Session.CompanyID,
		CampaignID: input.CampaignID,
		Status:     campaign_enum.Active,
		OR:         OR,
	})

	if err != nil {
		return
	}

	if len(companyAvg) == 0 {
		return
	}

	out.CompanyAVG = companyAvg[0].AVG
	return
}

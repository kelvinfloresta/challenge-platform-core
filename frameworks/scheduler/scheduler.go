package scheduler

import (
	"conformity-core/usecases/billing_case"
	"conformity-core/usecases/campaign_case"
	"time"

	"github.com/robfig/cron/v3"
)

const EVERY_DAY_AT_23_HOURS = "0 23 * * *"
const EVERY_DAY_1_AT_12_HOURS = "0 12 1 * *"

func Start() error {
	c := cron.New(cron.WithLocation(time.UTC))

	_, err := c.AddFunc(EVERY_DAY_AT_23_HOURS, func() {
		err := campaign_case.Singleton.NotifyDeadlineChallenge(time.Now())
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		return err
	}

	_, err = c.AddFunc(EVERY_DAY_AT_23_HOURS, func() {
		err := campaign_case.Singleton.NotifyNewChallenge(time.Now())
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		return err
	}

	_, err = c.AddFunc(EVERY_DAY_AT_23_HOURS, func() {
		err := campaign_case.Singleton.NotifyManagerDeadlineChallenge(time.Now())
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		return err
	}

	_, err = c.AddFunc(EVERY_DAY_1_AT_12_HOURS, func() {
		err := billing_case.Singleton.SendBillingReport(time.Now())
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		return err
	}

	go c.Start()

	return nil
}

package cronjob

import (
	"ferry/pkg/logger"

	"github.com/robfig/cron/v3"
)

func TestJob(c *cron.Cron) {
	id, err := c.AddFunc("1 * * * *", func() {

		logger.Info("Every hour on the one hour")
	})
	if err != nil {
		logger.Info(err)
		logger.Info("start error")
	} else {
		logger.Infof("Start Success; ID: %v \r\n", id)
	}
}

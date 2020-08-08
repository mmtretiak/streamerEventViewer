package jobs

import (
	"github.com/labstack/echo"
	"github.com/robfig/cron/v3"
	"streamerEventViewer/cmd/config"
	"streamerEventViewer/internal/jobs/view_updater"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/models"
)

func StartJobs(config config.Jobs, clipRepository models.ClipRepository, helixService helixService.Service, logger echo.Logger) {
	cronRunner := cron.New()
	viewUpdater := view_updater.New(clipRepository, helixService, logger)

	cronRunner.Schedule(config.ViewUpdaterJob.Schedule, viewUpdater)
}

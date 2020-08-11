package jobs

import (
	"github.com/labstack/echo"
	"github.com/robfig/cron/v3"
	"streamerEventViewer/cmd/config"
	"streamerEventViewer/internal/jobs/view_updater"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/models"
	"time"
)

func StartJobs(config config.Jobs, clipRepository models.ClipRepository, helixService helixService.Service, logger echo.Logger) {
	cronRunner := cron.New(cron.WithLocation(time.Local))
	viewUpdater := view_updater.New(clipRepository, helixService, logger)

	cronRunner.AddJob(config.ViewUpdaterJob.Schedule, viewUpdater)
	cronRunner.Start()
}

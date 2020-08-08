package view_updater

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/nicklaw5/helix"
	"github.com/robfig/cron/v3"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/models"
)

type viewUpdater struct {
	clipRepository models.ClipRepository
	helixService   helixService.Service
	logger         echo.Logger
}

// Run runs updater in goroutine
func New(clipRepository models.ClipRepository, helixService helixService.Service, logger echo.Logger) cron.Job {
	updater := viewUpdater{
		helixService:   helixService,
		clipRepository: clipRepository,
		logger:         logger,
	}

	return &updater
}

func (v *viewUpdater) Run() {
	ctx := context.Background()

	clips, err := v.getClips(ctx)
	if err != nil {
		v.logger.Error(err)
		return
	}

	v.updateClipViews(ctx, clips)
}

func (v *viewUpdater) getClips(ctx context.Context) ([]helix.Clip, error) {
	clips, err := v.clipRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all clips, reason: %v", err)
	}

	helixClient, err := v.helixService.NewAppClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create helix app client, reason: %v", err)
	}

	clipParams := &helix.ClipsParams{}

	// TODO note that twitch support not more than 100 clips IDs in request, so need to take care of this later
	for _, clip := range clips {
		clipParams.IDs = append(clipParams.IDs, clip.ExternalID)
	}

	resp, err := helixClient.GetClips(clipParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get clips info, reason: %v", err)
	}

	return resp.Data.Clips, nil
}

func (v *viewUpdater) updateClipViews(ctx context.Context, clips []helix.Clip) {
	// TODO investigate bulk update
	for _, clip := range clips {
		err := v.clipRepository.UpdateViewCountByExternalID(ctx, clip.ID, clip.ViewCount)
		if err != nil {
			v.logger.Errorf("failed to update clip with external id %s, reason: %v", clip.ID, err)
		}
	}
}

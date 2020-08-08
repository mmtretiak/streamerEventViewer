package clip

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/nicklaw5/helix"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/models"
	"streamerEventViewer/pkg/rbac"
)

type Service interface {
	SaveClip(c echo.Context, externalStreamerID string) error
	GetClipsForStreamer(c echo.Context, streamerID string) ([]models.Clip, error)
}

func New(helixService helixService.Service, rbacService rbac.RBACService, userSecretRepository models.UserSecretRepository,
	clipRepository models.ClipRepository, streamerRepository models.StreamerRepository, logger echo.Logger) Service {
	return &service{
		helixService:         helixService,
		rbac:                 rbacService,
		userSecretRepository: userSecretRepository,
		clipRepository:       clipRepository,
		streamerRepository:   streamerRepository,
		logger:               logger,
	}
}

type service struct {
	helixService         helixService.Service
	rbac                 rbac.RBACService
	userSecretRepository models.UserSecretRepository
	clipRepository       models.ClipRepository
	streamerRepository   models.StreamerRepository
	logger               echo.Logger
}

func (s *service) SaveClip(c echo.Context, externalStreamerID string) error {
	user := s.rbac.User(c)

	ctx := context.Background()

	streamer, err := s.streamerRepository.GetByExternalID(ctx, externalStreamerID)
	if err != nil {
		s.logger.Errorf("failed to get streamer by external ID %s, reason: %v", externalStreamerID, err)
		return err
	}

	secret, err := s.userSecretRepository.GetByUserID(ctx, user.ID)
	if err != nil {
		s.logger.Errorf("failed to get user secret by user ID %s, reason: %v", user.ID, err)
		return err
	}

	userHelixClient, err := s.helixService.NewUserClient(secret.AuthToken)
	if err != nil {
		s.logger.Errorf("failed to create user helix client for user ID %s, reason: %v", user.ID, err)
		return err
	}

	createClipParams := &helix.CreateClipParams{
		BroadcasterID: externalStreamerID,
		// TODO provide this in request, so user would have choice
		HasDelay: false,
	}

	resp, err := userHelixClient.CreateClip(createClipParams)
	if err != nil {
		s.logger.Errorf("failed to create clip for streamer ID %s, reason: %v", externalStreamerID, err)
		return err
	}

	clipInfo := resp.Data.ClipEditURLs[0]

	saveClipReq := saveClipReq{
		clipInfo:   clipInfo,
		userID:     user.ID,
		streamerID: streamer.ID,
	}

	err = s.saveClip(ctx, saveClipReq)
	if err != nil {
		s.logger.Errorf("failed to save clip external ID %s for streamer ID %s, reason: %v", clipInfo.ID, externalStreamerID, err)
		return err
	}

	return nil
}

func (s *service) GetClipsForStreamer(c echo.Context, streamerID string) ([]models.Clip, error) {
	user := s.rbac.User(c)

	ctx := context.Background()

	clips, err := s.clipRepository.GetByUserAndStreamerID(ctx, user.ID, streamerID)
	if err != nil {
		s.logger.Errorf("failed to get clip by user ID %s and streamer ID %s, reason: %v", user.ID, streamerID, err)
		return nil, err
	}

	return clips, err
}

type saveClipReq struct {
	clipInfo   helix.ClipEditURL
	userID     string
	streamerID string
}

func (s *service) saveClip(ctx context.Context, req saveClipReq) error {
	clip := models.Clip{
		ID:         uuid.New().String(),
		UserID:     req.userID,
		StreamerID: req.streamerID,
		ExternalID: req.clipInfo.ID,
		EditURL:    req.clipInfo.EditURL,
	}

	return s.clipRepository.Save(ctx, clip)
}

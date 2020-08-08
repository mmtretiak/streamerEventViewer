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
	clipRepository models.ClipRepository, streamerRepository models.StreamerRepository) Service {
	return &service{
		helixService:         helixService,
		rbac:                 rbacService,
		userSecretRepository: userSecretRepository,
		clipRepository:       clipRepository,
		streamerRepository:   streamerRepository,
	}
}

type service struct {
	helixService         helixService.Service
	rbac                 rbac.RBACService
	userSecretRepository models.UserSecretRepository
	clipRepository       models.ClipRepository
	streamerRepository   models.StreamerRepository
}

func (s *service) SaveClip(c echo.Context, externalStreamerID string) error {
	user := s.rbac.User(c)

	ctx := context.Background()

	streamer, err := s.streamerRepository.GetByExternalID(ctx, externalStreamerID)
	if err != nil {
		return err
	}

	secret, err := s.userSecretRepository.GetByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	userHelixClient, err := s.helixService.NewUserClient(secret.AuthToken)
	if err != nil {
		return err
	}

	createClipParams := &helix.CreateClipParams{
		BroadcasterID: externalStreamerID,
		// TODO provide this in request, so user would have choice
		HasDelay: false,
	}

	resp, err := userHelixClient.CreateClip(createClipParams)
	if err != nil {
		return err
	}

	clipInfo := resp.Data.ClipEditURLs[0]

	saveClipReq := saveClipReq{
		clipInfo:   clipInfo,
		userID:     user.ID,
		streamerID: streamer.ID,
	}

	return s.saveClip(ctx, saveClipReq)
}

func (s *service) GetClipsForStreamer(c echo.Context, streamerID string) ([]models.Clip, error) {
	user := s.rbac.User(c)

	ctx := context.Background()

	clips, err := s.clipRepository.GetByUserAndStreamerID(ctx, user.ID, streamerID)
	if err != nil {
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

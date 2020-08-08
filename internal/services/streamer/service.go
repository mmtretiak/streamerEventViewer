package streamer

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"streamerEventViewer/pkg/models"
	"streamerEventViewer/pkg/rbac"
)

type Service interface {
	SaveStreamer(c echo.Context, streamerName string) error
	GetStreamersForUser(echo.Context) ([]models.Streamer, error)
}

func New(rbacService rbac.RBACService, streamerRepository models.StreamerRepository, usersToStreamersRepository models.UsersToStreamerRepository) Service {
	return &service{
		rbac:                      rbacService,
		streamerRepository:        streamerRepository,
		usersToStreamerRepository: usersToStreamersRepository,
	}
}

type service struct {
	rbac                      rbac.RBACService
	streamerRepository        models.StreamerRepository
	usersToStreamerRepository models.UsersToStreamerRepository
}

func (s *service) SaveStreamer(c echo.Context, streamerName string) error {
	user := s.rbac.User(c)

	ctx := context.Background()

	var streamer models.Streamer
	var err error
	streamer, err = s.streamerRepository.GetByName(ctx, streamerName)
	if err != nil {
		streamer, err = s.createStreamer(ctx, streamerName)
		if err != nil {
			return err
		}
	}

	userToStreamer := models.UserToStreamer{
		StreamerID: streamer.ID,
		UserID:     user.ID,
	}

	isExists, err := s.usersToStreamerRepository.IsExist(ctx, userToStreamer)
	if err != nil {
		return err
	}

	if isExists {
		return nil
	}

	return s.usersToStreamerRepository.Save(ctx, userToStreamer)
}

func (s *service) createStreamer(ctx context.Context, streamerName string) (models.Streamer, error) {
	streamer := models.Streamer{
		ID:   uuid.New().String(),
		Name: streamerName,
	}

	if err := s.streamerRepository.Save(ctx, streamer); err != nil {
		return models.Streamer{}, err
	}

	return streamer, nil
}

func (s *service) GetStreamersForUser(c echo.Context) ([]models.Streamer, error) {
	user := s.rbac.User(c)

	ctx := context.Background()

	streamers, err := s.streamerRepository.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return streamers, nil
}

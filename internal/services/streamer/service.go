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
}

func New() Service {
	return &service{}
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

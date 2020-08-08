package streamer

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/nicklaw5/helix"
	"net/http"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/models"
	"streamerEventViewer/pkg/rbac"
)

type Service interface {
	SaveStreamer(c echo.Context, streamerName string) (int, error)
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
	helixService              helixService.Service
}

func (s *service) SaveStreamer(c echo.Context, streamerName string) (int, error) {
	searchResp, err := s.searchForStreamer(streamerName)
	if err != nil {
		return searchResp.status, err
	}

	user := s.rbac.User(c)

	ctx := context.Background()

	var streamer models.Streamer
	streamer, err = s.streamerRepository.GetByName(ctx, streamerName)
	if err != nil {
		streamer, err = s.createStreamer(ctx, searchResp.channel)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	userToStreamer := models.UserToStreamer{
		StreamerID: streamer.ID,
		UserID:     user.ID,
	}

	isExists, err := s.usersToStreamerRepository.IsExist(ctx, userToStreamer)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if isExists {
		return http.StatusInternalServerError, err
	}

	err = s.usersToStreamerRepository.Save(ctx, userToStreamer)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *service) createStreamer(ctx context.Context, channel helix.Channel) (models.Streamer, error) {
	streamer := models.Streamer{
		ID:         uuid.New().String(),
		Name:       channel.DisplayName,
		ExternalID: channel.ID,
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

type searchResponse struct {
	channel helix.Channel
	status  int
}

func (s *service) searchForStreamer(streamerName string) (searchResponse, error) {
	searchResp := searchResponse{}

	helixClient, err := s.helixService.NewAppClient()
	if err != nil {
		searchResp.status = http.StatusInternalServerError
		return searchResp, err
	}

	searchChannels := &helix.SearchChannelsParams{
		Channel: streamerName,
	}

	resp, err := helixClient.SearchChannels(searchChannels)
	if err != nil {
		searchResp.status = http.StatusBadRequest
		return searchResp, err
	}

	if len(resp.Data.Channels) == 0 {
		searchResp.status = http.StatusBadRequest
		return searchResp, errors.New(fmt.Sprintf("can not find streamer with name %s", streamerName))
	}

	// TODO make suggestion in better way, for example return them as separate field instead of error and parse on front-end
	if len(resp.Data.Channels) > 1 {
		errMsg := "please specify streamer name, suggestions %v"

		var suggestions []string
		for _, channel := range resp.Data.Channels {
			suggestions = append(suggestions, channel.DisplayName)
		}

		searchResp.status = http.StatusMultipleChoices

		return searchResp, errors.New(fmt.Sprintf(errMsg, suggestions))
	}

	searchResp.channel = resp.Data.Channels[0]
	return searchResp, nil
}

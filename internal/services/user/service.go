package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/nicklaw5/helix"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/models"
)

// TODO generate this for each user and put into cache(redis for example)
const state = "test-state"

type Service interface {
}

type service struct {
	helixService         helixService.Service
	userRepository       models.UserRepository
	userSecretRepository models.UserSecretRepository
}

func New() Service {
	return &service{}
}

func (s *service) Login(c echo.Context) error {
	helixClient, err := s.helixService.NewAppClient()
	if err != nil {
		return err
	}

	url := helixClient.GetAuthorizationURL(state, true)

	return c.Redirect(200, url)
}

func (s *service) Redirect(c echo.Context) error {
	authCode := c.FormValue("code")

	helixAppClient, err := s.helixService.NewAppClient()
	if err != nil {
		return err
	}

	userAccessTokenResp, err := helixAppClient.GetUserAccessToken(authCode)
	if err != nil {
		return err
	}

	accessToken := userAccessTokenResp.Data.AccessToken

	helixUserClient, err := s.helixService.NewUserClient(accessToken)
	if err != nil {
		return err
	}

	userInfoResp, err := helixUserClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		return err
	}

	users := userInfoResp.Data.Users

	if len(users) != 1 {
		return errors.New(fmt.Sprintf("received %v users instead of 1", len(users)))
	}

	user := convertUser(users[0])
	secret := convertSecret(user.ID, userAccessTokenResp.Data)

	// TODO should be context with deadlines
	ctx := context.Background()

	if err := s.userRepository.Save(ctx, user); err != nil {
		return err
	}

	if err := s.userSecretRepository.Save(ctx, secret); err != nil {
		// TODO should be deletion of user
		return err
	}

	return nil
}

func convertSecret(userID string, credentials helix.UserAccessCredentials) models.UserSecret {
	convertedSecret := models.UserSecret{
		ID:        uuid.New().String(),
		UserID:    userID,
		AuthToken: credentials.AccessToken,
		Scopes:    credentials.Scopes,
	}

	return convertedSecret
}

func convertUser(user helix.User) models.User {
	convertedUser := models.User{
		ID:           uuid.New().String(),
		Email:        user.Email,
		Name:         user.Login,
		ThumbnailURL: user.ProfileImageURL,
	}

	return convertedUser
}

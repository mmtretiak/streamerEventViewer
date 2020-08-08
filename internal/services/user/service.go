package user

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
)

// TODO generate this for each user and put into cache(redis for example)
const state = "test-state"

type Service interface {
	Login(c echo.Context) error
	Redirect(c echo.Context, authCode string) (models.Token, error)
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(models.User) (string, error)
}

type service struct {
	helixService         helixService.Service
	userRepository       models.UserRepository
	userSecretRepository models.UserSecretRepository
	tokenGenerator       TokenGenerator
	logger               echo.Logger
}

func New(helixService helixService.Service, userRepo models.UserRepository, userSecretRepo models.UserSecretRepository,
	generator TokenGenerator, logger echo.Logger) Service {
	return &service{
		helixService:         helixService,
		userSecretRepository: userSecretRepo,
		userRepository:       userRepo,
		tokenGenerator:       generator,
		logger:               logger,
	}
}

func (s *service) Login(c echo.Context) error {
	helixClient, err := s.helixService.NewAppClient()
	if err != nil {
		s.logger.Errorf("failed to create app client, reason: %v", err)
		return err
	}

	url := helixClient.GetAuthorizationURL(state, true)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *service) Redirect(c echo.Context, authCode string) (models.Token, error) {
	helixAppClient, err := s.helixService.NewAppClient()
	if err != nil {
		s.logger.Errorf("failed to create app client, reason: %v", err)
		return models.Token{}, err
	}

	userAccessTokenResp, err := helixAppClient.GetUserAccessToken(authCode)
	if err != nil {
		s.logger.Errorf("failed to get user access token, reason: %v", err)
		return models.Token{}, err
	}

	accessToken := userAccessTokenResp.Data.AccessToken

	helixUserClient, err := s.helixService.NewUserClient(accessToken)
	if err != nil {
		s.logger.Errorf("failed to create user client, reason: %v", err)
		return models.Token{}, err
	}

	userInfoResp, err := helixUserClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		s.logger.Errorf("failed to get users info, reason: %v", err)
		return models.Token{}, err
	}

	users := userInfoResp.Data.Users

	if len(users) != 1 {
		err = errors.New(fmt.Sprintf("received %v users instead of 1", len(users)))
		s.logger.Error(err)

		return models.Token{}, err
	}

	user := convertUser(users[0])
	secret := convertSecret(user.ID, userAccessTokenResp.Data)

	// TODO should be context with deadlines
	ctx := context.Background()

	token, err := s.tokenGenerator.GenerateToken(user)
	if err != nil {
		s.logger.Errorf("failed to generate token for user ID %s, reason: %v", user.ID, err)

		return models.Token{}, err
	}

	if err := s.userRepository.Save(ctx, user); err != nil {
		s.logger.Errorf("failed to save user ID %s, reason: %v", user.ID, err)
		return models.Token{}, err
	}

	if err := s.userSecretRepository.Save(ctx, secret); err != nil {
		// TODO should be deletion of user
		s.logger.Errorf("failed to save user secret, user ID %s, secret ID %s, reason: %v", user.ID, secret.ID, err)

		return models.Token{}, err
	}

	return models.Token{Token: token}, nil
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

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
	Redirect(c echo.Context) (models.Token, error)
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
}

func New(helixService helixService.Service, userRepo models.UserRepository, userSecretRepo models.UserSecretRepository,
	generator TokenGenerator) Service {
	return &service{
		helixService:         helixService,
		userSecretRepository: userSecretRepo,
		userRepository:       userRepo,
		tokenGenerator:       generator,
	}
}

func (s *service) Login(c echo.Context) error {
	helixClient, err := s.helixService.NewAppClient()
	if err != nil {
		return err
	}

	url := helixClient.GetAuthorizationURL(state, true)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *service) Redirect(c echo.Context) (models.Token, error) {
	authCode := c.FormValue("code")

	helixAppClient, err := s.helixService.NewAppClient()
	if err != nil {
		return models.Token{}, err
	}

	userAccessTokenResp, err := helixAppClient.GetUserAccessToken(authCode)
	if err != nil {
		return models.Token{}, err
	}

	accessToken := userAccessTokenResp.Data.AccessToken

	helixUserClient, err := s.helixService.NewUserClient(accessToken)
	if err != nil {
		return models.Token{}, err
	}

	userInfoResp, err := helixUserClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		return models.Token{}, err
	}

	users := userInfoResp.Data.Users

	if len(users) != 1 {
		return models.Token{}, errors.New(fmt.Sprintf("received %v users instead of 1", len(users)))
	}

	user := convertUser(users[0])
	secret := convertSecret(user.ID, userAccessTokenResp.Data)

	// TODO should be context with deadlines
	ctx := context.Background()

	token, err := s.tokenGenerator.GenerateToken(user)
	if err != nil {
		return models.Token{}, err
	}

	if err := s.userRepository.Save(ctx, user); err != nil {
		return models.Token{}, err
	}

	if err := s.userSecretRepository.Save(ctx, secret); err != nil {
		// TODO should be deletion of user
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

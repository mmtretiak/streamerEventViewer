package helix

import (
	"github.com/nicklaw5/helix"
	"streamerEventViewer/cmd/config"
)

func New(config config.OauthConfig) Service {
	return &helixService{
		config: config,
	}
}

type Service interface {
	NewAppClient() (*helix.Client, error)
	NewUserClient(userAccessToken string) (*helix.Client, error)
}

type helixService struct {
	config config.OauthConfig
}

func (h *helixService) NewAppClient() (*helix.Client, error) {
	options := &helix.Options{
		ClientID:     h.config.ClientID,
		ClientSecret: h.config.ClientSecret,
		RedirectURI:  h.config.RedirectURI,
		Scopes:       h.config.Scopes,
	}

	return helix.NewClient(options)
}

func (h *helixService) NewUserClient(userAccessToken string) (*helix.Client, error) {
	options := &helix.Options{
		ClientID:        h.config.ClientID,
		UserAccessToken: userAccessToken,
	}

	return helix.NewClient(options)
}

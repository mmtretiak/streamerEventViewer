package helix

import (
	"github.com/mmtretiak/helix"
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

	client, err := helix.NewClient(options)
	if err != nil {
		return nil, err
	}

	// TODO this not needed for all calls through this app client, so should move this creation into separate method of service
	resp, err := client.GetAppAccessToken()
	if err != nil {
		return nil, err
	}

	client.SetAppAccessToken(resp.Data.AccessToken)

	return client, nil
}

func (h *helixService) NewUserClient(userAccessToken string) (*helix.Client, error) {
	options := &helix.Options{
		ClientID:        h.config.ClientID,
		UserAccessToken: userAccessToken,
	}

	return helix.NewClient(options)
}

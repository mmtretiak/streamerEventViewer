package main

import (
	"fmt"
	"github.com/labstack/echo"
	"streamerEventViewer/cmd/config"
	"streamerEventViewer/internal/services/user"
	"streamerEventViewer/internal/services/user/transport"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/postgres"
	userRepo "streamerEventViewer/pkg/postgres/repositories/user"
	userSecretRepo "streamerEventViewer/pkg/postgres/repositories/user_secret"
	"streamerEventViewer/pkg/session"
)

func main() {
	config := config.New()

	db, err := postgres.New(config.DB)
	if err != nil {
		panic(err)
	}

	userRepository := userRepo.New(db)
	userSecretRepository := userSecretRepo.New(db)
	helixService := helixService.New(config.OauthConfig)

	jwt, err := session.New(config.JWTConfig, 0)
	if err != nil {
		panic(err)
	}

	userService := user.New(helixService, userRepository, userSecretRepository, jwt)

	e := echo.New()
	transport.NewHTTP(userService, e.Group(""))

	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	err = e.Start(address)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"github.com/labstack/echo"
	"streamerEventViewer/cmd/config"
	"streamerEventViewer/internal/services/clip"
	clipTransport "streamerEventViewer/internal/services/clip/transport"
	"streamerEventViewer/internal/services/streamer"
	streamerTransport "streamerEventViewer/internal/services/streamer/transport"
	"streamerEventViewer/internal/services/user"
	userTransport "streamerEventViewer/internal/services/user/transport"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/middleware/auth"
	"streamerEventViewer/pkg/postgres"
	clipRepo "streamerEventViewer/pkg/postgres/repositories/clip"
	streamerRepo "streamerEventViewer/pkg/postgres/repositories/streamer"
	userRepo "streamerEventViewer/pkg/postgres/repositories/user"
	userSecretRepo "streamerEventViewer/pkg/postgres/repositories/user_secret"
	"streamerEventViewer/pkg/postgres/repositories/user_to_streamer"
	"streamerEventViewer/pkg/rbac"
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
	streamerRepository := streamerRepo.New(db)
	userToStreamerRepository := user_to_streamer.New(db)
	clipRepository := clipRepo.New(db)

	rbacService := rbac.New()

	helixService := helixService.New(config.OauthConfig)

	jwt, err := session.New(config.JWTConfig, 0)
	if err != nil {
		panic(err)
	}

	authMiddleware := auth.Middleware(jwt)

	userService := user.New(helixService, userRepository, userSecretRepository, jwt)
	streamerService := streamer.New(rbacService, streamerRepository, userToStreamerRepository)
	clipService := clip.New(helixService, rbacService, userSecretRepository, clipRepository, streamerRepository)

	e := echo.New()
	eGroup := e.Group("")

	userTransport.NewHTTP(userService, eGroup)
	streamerTransport.NewHTTP(streamerService, eGroup, authMiddleware)
	clipTransport.NewHTTP(clipService, eGroup, authMiddleware)

	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	err = e.Start(address)
	if err != nil {
		panic(err)
	}
}

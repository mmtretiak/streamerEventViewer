package main

import (
	"fmt"
	"github.com/labstack/echo"
	"os"
	"streamerEventViewer/cmd/config"
	"streamerEventViewer/internal/jobs"
	"streamerEventViewer/internal/services/clip"
	clipTransport "streamerEventViewer/internal/services/clip/transport"
	"streamerEventViewer/internal/services/streamer"
	streamerTransport "streamerEventViewer/internal/services/streamer/transport"
	"streamerEventViewer/internal/services/user"
	userTransport "streamerEventViewer/internal/services/user/transport"
	helixService "streamerEventViewer/pkg/helix"
	"streamerEventViewer/pkg/middleware/auth"
	"streamerEventViewer/pkg/middleware/secure"
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

	e := echo.New()
	e.Use(secure.CORS(), secure.Headers())

	logger := e.Logger

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

	userService := user.New(helixService, userRepository, userSecretRepository, jwt, logger)
	streamerService := streamer.New(rbacService, streamerRepository, userToStreamerRepository, helixService, logger)
	clipService := clip.New(helixService, rbacService, userSecretRepository, clipRepository, streamerRepository, logger)

	eGroup := e.Group("")

	userTransport.NewHTTP(userService, eGroup)
	streamerTransport.NewHTTP(streamerService, eGroup, authMiddleware)
	clipTransport.NewHTTP(clipService, eGroup, authMiddleware)

	address := GetAddress(config.Server)

	go jobs.StartJobs(config.Jobs, clipRepository, helixService, logger)

	echo.NotFoundHandler = func(c echo.Context) error {
		return c.File("cmd/dist/index.html")
	}

	e.Static("/", "cmd/dist")

	err = e.Start(address)
	if err != nil {
		panic(err)
	}
}

func GetAddress(config config.Server) string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		return fmt.Sprintf("%s:%d", config.Host, config.Port)
	}

	return ":" + port
}

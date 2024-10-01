package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"gitlab.figvam.ru/figvam/tgapi/internal/config"
	"gitlab.figvam.ru/figvam/tgapi/internal/pg"
	"gitlab.figvam.ru/figvam/tgapi/internal/secret"
	"gitlab.figvam.ru/figvam/tgapi/internal/server"
	"gitlab.figvam.ru/figvam/tgapi/internal/tgap"
	"gitlab.figvam.ru/figvam/tgapi/pkg"
)

func Run() {

	cfg, err := config.GetAppConfig()
	if err != nil || cfg == nil {
		panic(err)
	}

	log.Logger = pkg.NewLogger(cfg.LogCfg)

	logger := log.Logger

	log.Info().Str("comp:", "main").Msg("log initiated")

	mainCtx, cancelMainCtx := context.WithCancel(context.Background())

	db := pg.New(mainCtx, cfg, logger)

	vt := secret.NewVault(mainCtx, cfg, logger)

	handler := server.NewHandler(logger)

	router := handler.SetRoutes()
	log.Info().Str("comp:", "app/set routes").Msg("Routes set")

	go func() {
		err := router.Start(":" + cfg.ServicePort)
		if err != nil {
			log.Fatal().Err(err).Str("Server", "Start").Msg("Error while starting server")
		}
	}()

	go func() {
		if err := tgap.NewClient(mainCtx, cfg, logger, db, vt); err != nil {
			panic(err)
		}

	}()

	//Wait kill signal
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)
	<-killSignal
	cancelMainCtx()
	logger.Info().Str("comp:", "main").Msg("Graceful shutdown. This can take a while...")
	err = secret.WriteSecret(vt, cfg, logger)
	if err != nil {
		logger.Err(err).Str("Graceful", "Write Secret").Msg("error while writing secret")
	} else {
		logger.Info().Str("Graceful", "Write Secret").Msg("Successful")
	}

}

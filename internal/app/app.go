package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"gitlab.figvam.ru/figvam/tgapi/internal/config"
	"gitlab.figvam.ru/figvam/tgapi/internal/pg"
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

	handler := server.NewHandler(logger)

	router := handler.SetRoutes()
	log.Info().Str("comp:", "app/set routes").Msg("Routes set")
	go func() {

		err := router.Start(":" + cfg.ServicePort)
		if err != nil {
			log.Fatal().Err(err).Str("Server", "Start").Msg("Error while starting server")
		}
	}()
	logger.Info().Str("comp:", "tgap").Any("ID:=", cfg.ID).Msg("Debug LOG")
	go func() {

		if err := tgap.New(cfg, logger, db).NewClient(mainCtx); err != nil {
			log.Fatal().Err(err).Str("Tgap", "Start").Msg("Error while starting tgap client")
		}
	}()

	//Wait kill signal
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

	<-killSignal
	logger.Info().Str("comp:", "main").Msg("Graceful shutdown. This can take a while...")

	cancelMainCtx()

}

package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tgapiV2/internal/config"
	"tgapiV2/internal/pg"
	"tgapiV2/internal/secret"
	"tgapiV2/internal/server"
	"time"

	"tgapiV2/internal/tgap"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run(
	cfg *config.Appconfig,
	logger zerolog.Logger) error {

	mainCtx, cancelMainCtx := context.WithCancel(context.Background())

	db := pg.New(mainCtx, cfg, logger)

	vt := secret.NewVault(mainCtx, cfg, logger)

	handler := server.NewHandler(logger)

	router := handler.SetRoutes()
	log.Info().Str("comp:", "app/set routes").Msg("Routes set")
	srv := http.Server{
		Addr:              ":" + cfg.ServicePort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	//Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Str("comp:", "app/listenAndServe").Msg("error while start server")
			return
		}
	}()

	go func() {
		if client := tgap.NewClient(mainCtx, cfg, logger, db, vt); client != nil {
			panic(client)
		}

	}()

	//Wait kill signal
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)
	<-killSignal
	cancelMainCtx()
	logger.Info().Str("comp:", "main").Msg("Graceful shutdown. This can take a while...")
	err := secret.WriteSecret(vt, cfg, logger)
	if err != nil {
		logger.Err(err).Str("Graceful", "Write Secret").Msg("error while writing secret")
	} else {
		logger.Info().Str("Graceful", "Write Secret").Msg("Successful")
	}

	return nil

}

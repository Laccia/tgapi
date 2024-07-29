package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"myapp/app/configs"
	"myapp/delivery"
	"myapp/hide"
	"myapp/logic"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/go-faster/errors"
	"github.com/gotd/td/bin"
	"github.com/gotd/td/examples"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"
	vault "github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {

	config := vault.DefaultConfig()

	config.Address = "http://127.0.0.1:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken("hvs.xOkDyN9Ygg4ApFc5wDRynbZs")
	hide.ReadSecret(client)
	time.Sleep(1 * time.Second)
	log, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel), zap.AddStacktrace(zapcore.FatalLevel))
	defer func() { _ = log.Sync() }()

	cfg := configs.New()
	logic.New(cfg)
	pg := delivery.New(cfg)

	d := tg.NewUpdateDispatcher()
	s := tg.NewServerDispatcher(func(ctx context.Context, b *bin.Buffer) (a bin.Encoder, e error) {
		fmt.Println("maybe!!!!!!!!!!!!!!!!!!!!!!!!!")
		return a, e
	})
	b := bin.Buffer{}

	gaps := updates.New(updates.Config{
		Handler: d,
		Logger:  log.Named("gaps"),
	})

	// Authentication flow handles authentication process, like prompting for code and 2FA password.
	flow := auth.NewFlow(examples.Terminal{}, auth.SendCodeOptions{})
	// Initializing client from environment.
	// Available environment variables:
	// APP_ID: ""
	// APP_HASH:       app_hash of Telegram app.
	// SESSION_FILE:   path to session file
	// SESSION_DIR:    path to session directory, if SESSION_FILE is not set
	tgap, err := telegram.ClientFromEnvironment(telegram.Options{
		Logger:        log,
		UpdateHandler: gaps,
		Middlewares: []telegram.Middleware{
			updhook.UpdateHook(gaps.Handle),
		},
	})
	s.Handle(ctx, &b)

	// Setup message update handlers.
	d.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		newmsg, err := json.MarshalIndent(update.Message, "", "")
		if err != nil {
			log.Info("cant marshal")
		}

		file := os.Getenv("TEMPLATE_FILE")
		tfile, err := os.ReadFile(file)
		if err != nil {
			log.Info("cant read file")
		}

		msg := make(map[string]interface{})
		template := make(map[string]interface{})
		templ := json.Unmarshal(tfile, &template)
		if templ != nil {
			log.Info("cant marshall")
		}

		s := json.Unmarshal(newmsg, &msg)
		if s != nil {
			log.Info("cant unmarshall")
		}

		id := msg["PeerID"]
		tempID := template["PeerID"]
		tar := map[string]interface{}{
			"PeerID": id,
		}

		ter := map[string]interface{}{
			"PeerID": tempID,
		}
		// b, err := json.Marshal(&s)
		// if err != nil {
		// 	fmt.Println("cant marshal msg for pg")
		// }
		// var md models.Pgmsg
		// md.Msg = update.Message.String()
		// message := update.Message.String()
		equal := reflect.DeepEqual(tar, ter)
		if equal {
			err = pg.Logic.AddMsg(ctx, newmsg)
			if err != nil {
				fmt.Println(err)
			}

			log.Info(update.Message.String())
		} else {
			log.Info("Trash")
		}

		fmt.Println(id)
		ts := os.WriteFile("tmp.json", newmsg, 0644)
		if ts != nil {
			log.Info("cant write file")
		}

		return nil

	})
	// a, e := s.Handle(ctx, &bin.Buffer{})
	// if e != nil {
	// 	fmt.Println("YA EBAL")
	// }
	// fmt.Println(a)
	s.OnChannelsGetMessages(func(ctx context.Context, request *tg.ChannelsGetMessagesRequest) (a tg.MessagesMessagesClass, e error) {
		s := request.String()
		nn, bb := request.GetChannelAsNotEmpty()
		request.Channel.AsNotEmpty()
		fmt.Println(nn, bb)
		e = pg.Logic.AddHistory(ctx, s)

		if err != nil {
			log.Info("Cant add history")
		}
		return
	})
	s.OnMessagesGetHistory(func(ctx context.Context, request *tg.MessagesGetHistoryRequest) (a tg.MessagesMessagesClass, e error) {
		inp := request.GetPeer()
		request.Peer = inp
		request.Limit = 100
		request.GetLimit()
		c := request.String()
		fmt.Println(inp, c)

		return a, e

	})

	d.OnPeerHistoryTTL(func(ctx context.Context, e tg.Entities, update *tg.UpdatePeerHistoryTTL) error {

		err = pg.Logic.AddHistory(ctx, update.String())
		if err != nil {
			log.Info("Cant add history")
		}
		return err
	})

	return tgap.Run(ctx, func(ctx context.Context) error {
		// Perform auth if no session is available.
		if err := tgap.Auth().IfNecessary(ctx, flow); err != nil {
			return errors.Wrap(err, "auth")
		}

		// Fetch user info.
		user, err := tgap.Self(ctx)
		if err != nil {
			return errors.Wrap(err, "call self")
		}

		return gaps.Run(ctx, tgap.API(), user.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				hide.WriteSecret(client)
				log.Info("Gaps started")

			},
		})
	})

}

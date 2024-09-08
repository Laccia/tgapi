package secret

import (
	"encoding/json"
	"os"
	"tgapiV2/internal/config"

	vault "github.com/hashicorp/vault/api"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

func NewVault(ctx context.Context, cfg *config.Appconfig, log zerolog.Logger) *vault.Client {

	config := vault.DefaultConfig()

	config.Address = cfg.Host

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatal().Err(err).Str("vault", "client").Msg("error while create new client")
	}

	client.SetToken(cfg.Token)

	err = ReadSecret(client, cfg, log)
	if err != nil {
		log.Err(err).Str("vault", "read/secret").Msg("error while read secret from vault. Session file will create from auth")
	}

	return client
}

func ReadSecret(client *vault.Client, cfg *config.Appconfig, log zerolog.Logger) error {

	secret, err := client.KVv2(cfg.MountPath).Get(context.Background(), cfg.ReadPath)
	if err != nil {
		return err
	}

	file, _ := json.Marshal(secret.Data)

	err = os.WriteFile(cfg.File, file, 0777)
	if err != nil {
		return err
	}

	log.Info().Str("vault", "read/secret").Msg("secret apply")
	return nil
}

func WriteSecret(client *vault.Client, cfg *config.Appconfig, log zerolog.Logger) error {

	previous, err := os.ReadFile(cfg.File)
	if err != nil {
		return err
	}

	secretData := map[string]interface{}{}
	// json to map
	err = json.Unmarshal(previous, &secretData)

	if err != nil {
		return err
	}

	// Write a secret
	_, err = client.KVv2(cfg.MountPath).Put(context.Background(), cfg.WritePath, secretData)
	if err != nil {
		return err
	}
	return nil
}

package config

import (
	"fmt"

	"github.com/spf13/viper"
	"gitlab.figvam.ru/figvam/tgapi/pkg"
)

type Appconfig struct {
	Dbconfig
	TGconfig
	LogCfg      pkg.Config
	ServicePort string
}

type Dbconfig struct {
	URL string
}

type TGconfig struct {
	ID    int
	Hash  string
	File  string
	Dir   string
	Templ string
	Phone string
	Chats []int64
}

type AppConfigParseStruct struct {
	Port string `mapstructure:"PORT"`

	//DB config

	PGHost string `mapstructure:"PG_HOST"`
	PGName string `mapstructure:"PG_NAME"`
	PGPass string `mapstructure:"PG_PASS"`
	PGAddr string `mapstructure:"PG_ADDR"`
	PGPort string `mapstructure:"PG_PORT"`
	PGBase string `mapstructure:"PH_BASE"`

	//TG config
	APPID           int     `mapstructure:"APP_ID"`
	APPHash         string  `mapstructure:"APP_HASH"`
	SessionFile     string  `mapstructure:"SESSION_FILE"`
	SessionDir      string  `mapstructure:"SESSION_DIR"`
	SessionTemplate string  `mapstructure:"TG_SESSION_TEMPLATE"`
	PhoneNumber     string  `mapstructure:"PHONE_NUMBER"`
	ChatID          []int64 `mapstructure:"CHAT_ID"`

	//Logger settings
	// Logger settings
	LogLevel    pkg.LogLevel       `mapstructure:"LEVEL"`
	Format      pkg.LogFormat      `mapstructure:"FORMAT"`
	Destination pkg.LogDestination `mapstructure:"DESTINATION"`
}

func GetAppConfig() (*Appconfig, error) {
	viper.SetDefault("TG_SESSION_FILE", "utils/session.json")
	viper.SetDefault("TG_SESSION_DIR", "utils")
	viper.SetDefault("TG_SESSION_TEMPLATE", "/app/files/template.json")
	viper.SetDefault("PG_HOST", "postgres://postgres:password@localhost:5432/tg")
	viper.SetDefault("PORT", "8000")

	viper.SetDefault("LEVEL", "info")
	viper.SetDefault("FORMAT", 0)
	viper.SetDefault("DESTINATION", 3)
	viper.SetDefault("PHONE_NUMBER", "+79001411695")

	cfgParse := &AppConfigParseStruct{}
	cfg := &Appconfig{}
	if err := pkg.ParseConfig(cfgParse); err != nil {
		return cfg, fmt.Errorf("get config from environment variable: %w", err)
	}

	cfg.ServicePort = cfgParse.Port

	cfg.Dbconfig = Dbconfig{
		URL: cfgParse.PGHost,
	}

	cfg.TGconfig = TGconfig{
		ID:    cfgParse.APPID,
		Hash:  cfgParse.APPHash,
		File:  cfgParse.SessionFile,
		Dir:   cfgParse.SessionDir,
		Templ: cfgParse.SessionTemplate,
		Chats: cfgParse.ChatID,
		Phone: cfgParse.PhoneNumber,
	}

	cfg.LogCfg = pkg.Config{
		Level:       cfgParse.LogLevel,
		Destination: cfg.LogCfg.Destination,
		Format:      cfg.LogCfg.Format,
	}

	return cfg, nil
}

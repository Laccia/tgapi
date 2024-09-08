package config

import (
	"fmt"
	"tgapiV2/pkg"

	"github.com/spf13/viper"
)

type Appconfig struct {
	Dbconfig
	TGconfig
	VTconfig
	LogCfg      pkg.Config
	ServicePort string
}

type Dbconfig struct {
	URL  string
	Name string
	Pass string
	Addr string
	Port string
	Base string
}

type TGconfig struct {
	ID    string
	Hash  string
	File  string
	Dir   string
	Templ string
	Phone string
	Chats []int64
}

type VTconfig struct {
	Host      string
	Token     string
	MountPath string
	WritePath string
	ReadPath  string
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
	APPID           string  `mapstructure:"TG_APP_ID"`
	APPHash         string  `mapstructure:"TG_APP_HASH"`
	SessionFile     string  `mapstructure:"SESSION_FILE"`
	SessionDir      string  `mapstructure:"SESSION_DIR"`
	SessionTemplate string  `mapstructure:"TG_SESSION_TEMPLATE"`
	PhoneNumber     string  `mapstructure:"PHONE_NUMBER"`
	ChatID          []int64 `mapstructure:"CHAT_ID"`

	//VT config
	VaultHost      string `mapstructure:"VT_HOST"`
	VaultToken     string `mapstructure:"VT_TOKEN"`
	VaultMountPath string `mapstructure:"VT_MOUNT_PATH"`
	VaultWritePath string `mapstructure:"VT_WRITE_PATH"`
	VaultReadPath  string `mapstructure:"VT_READ_PATH"`

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

	viper.SetDefault("LEVEL", "info")
	viper.SetDefault("FORMAT", 0)
	viper.SetDefault("DESTINATION", 1)
	viper.SetDefault("PHONE_NUMBER", "+79001411695")

	cfgParse := &AppConfigParseStruct{}
	cfg := &Appconfig{}
	if err := pkg.ParseConfig(cfgParse); err != nil {
		return cfg, fmt.Errorf("get config from environment variable: %w", err)
	}

	cfg.ServicePort = cfgParse.Port

	cfg.Dbconfig = Dbconfig{
		URL:  cfgParse.PGHost,
		Name: cfgParse.PGName,
		Pass: cfgParse.PGPass,
		Addr: cfgParse.PGAddr,
		Port: cfgParse.PGPort,
		Base: cfgParse.PGBase,
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

	cfg.VTconfig = VTconfig{
		Host:      cfgParse.VaultHost,
		Token:     cfgParse.VaultToken,
		MountPath: cfgParse.VaultMountPath,
		WritePath: cfgParse.VaultWritePath,
		ReadPath:  cfgParse.VaultReadPath,
	}

	cfg.LogCfg = pkg.Config{
		Level:       cfgParse.LogLevel,
		Destination: cfg.LogCfg.Destination,
		Format:      cfg.LogCfg.Format,
	}

	return cfg, nil
}

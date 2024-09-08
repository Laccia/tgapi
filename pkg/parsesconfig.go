package pkg

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func ParseConfig(cfg interface{}) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	viper.AddConfigPath(workingDir)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// if config cannot be read, use env variables
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return err
		}
	}

	// Viper does not have a func that binds all the env variables
	// That's the reason why we get all the structure tags from config.ApplicationConfig
	envKeysMap := map[string]interface{}{}
	// We can use dummy structure, we only want to parse mapstructure tags
	if err := mapstructure.Decode(cfg, &envKeysMap); err != nil {
		return err
	}

	// All the parsed keys are being bind to env variables via Viper
	for key := range envKeysMap {
		if err := viper.BindEnv(key); err != nil {
			return err
		}
	}

	// We must decode all the keys in the structure, because viper's feature with auto structures' keys decode has failed
	// Viper's developers comment:
	// 		The feature is now disabled by default and can be enabled using the viper_bind_struct build tag.
	// 		It's also considered experimental at this point, so breaking changes may be introduced in the future.
	// LINK: https://github.com/spf13/viper/releases/tag/v1.18.2
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	validate := validator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("mapstructure"), ",", 2)[0]
		return name
	})

	if err := validate.Struct(cfg); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return errors.New("invalid value for validation")
		}

		var errs error
		for _, err := range err.(validator.ValidationErrors) {
			errs = errors.Join(errs, err)
		}
		return errs
	}

	return nil
}

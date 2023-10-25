package config

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// MetaConfig holds configuration for environment name and service name
type MetaConfig struct {
	Environment string
	ServiceName string `mapstructure:"service_name"`
}

// Load config from file then from environment variables
func Load(fs embed.FS, config interface{}) error {
	configType := "json"
	viper.SetConfigType(configType)

	configName := "config"
	if isTesting() {
		configName = "config_test"
	}
	viper.SetConfigName(configName)

	configFile := configName + "." + configType

	contents, err := fs.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("unable to read config from %s: %w", configFile, err)
	}

	if err := viper.ReadConfig(bytes.NewReader(contents)); err != nil {
		return fmt.Errorf("unable to read config from %s: %w", configFile, err)
	}

	return unmarshalConfig(config)
}

func unmarshalConfig(config interface{}) error {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read the config file again and consider environment variables at the same time
	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf(
			"unable to unmarshal config at %s: %w",
			viper.ConfigFileUsed(),
			err,
		)
	}

	// Set the environment to be "test" if tests are being run.
	meta := getMetaConfig(config)
	if meta == nil {
		return errors.New("meta config not available")
	}

	if isTesting() {
		meta.Environment = "test"
	} else {
		env := os.Getenv("ENV")
		if env != "" {
			meta.Environment = env
		}
	}

	return nil
}

func isTesting() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.v=") {
			return true
		}
	}
	return false || os.Getenv("ENV") == "test"
}

func getMetaConfig(config interface{}) *MetaConfig {
	v := reflect.ValueOf(config).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typ := field.Type()
		if typ.AssignableTo(reflect.TypeOf(MetaConfig{})) {
			return field.Addr().Interface().(*MetaConfig)
		}
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanAddr() {
			if meta := getMetaConfig(field.Addr().Interface()); meta != nil {
				return meta
			}
		}
	}
	return nil
}

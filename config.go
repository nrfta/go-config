package config

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"strings"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/neighborly/gtoolbox/errors"
	"github.com/spf13/viper"
)

type MetaConfig struct {
	Environment string
	ServiceName string `mapstructure:"service_name"`
}

func Load(box packr.Box, config interface{}) errors.Error {
	configType := "json"
	viper.SetConfigType(configType)

	configName := "config"
	if IsTesting() {
		configName = "config_test"
	}
	viper.SetConfigName(configName)

	configFile := configName + "." + configType
	contents, err := box.Find(configFile)
	if err != nil {
		return errors.Newf("unable to read config from %s", configFile).WithCause(err)
	}

	if err := viper.ReadConfig(bytes.NewReader(contents)); err != nil {
		return errors.Newf("unable to read config from %s", configFile).WithCause(err)
	}

	return unmarshalConfig(config)
}


func unmarshalConfig(config interface{}) errors.Error {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read the config file again and consider environment variables at the same time
	if err := viper.Unmarshal(config); err != nil {
		return errors.Newf("unable to unmarshal config at %s", viper.ConfigFileUsed()).WithCause(err)
	}

	// Set the environment to be "test" if tests are being run.
	meta := GetMetaConfig(config)
	if meta == nil {
		return errors.New("meta config not available")
	}

	if IsTesting() {
		meta.Environment = "test"
	} else {
		env := os.Getenv("ENV")
		if env != "" {
			meta.Environment = env
		}
	}

	return nil
}

func IsTesting() bool {
	return flag.Lookup("test.v") != nil || os.Getenv("ENV") == "test"
}

func GetMetaConfig(config interface{}) *MetaConfig {
	v := reflect.ValueOf(config).Elem()
	fmt.Println(config);
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
			if meta := GetMetaConfig(field.Addr().Interface()); meta != nil {
				return meta
			}
		}
	}
	return nil
}

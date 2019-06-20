package config

import (
	"bytes"
	"flag"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"

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

// Deprecated: Use Load(packr.Box, interface{}) to load config from the executable; override with environment variables.
func LoadConfig(config interface{}) errors.Error {
	viper.SetConfigType("json")

	execPath, err := GetExecutablePath()
	if err != nil {
		return err
	}
	viper.AddConfigPath(GetConfigPath(execPath))

	callerPath, err := GetCallerFilePath(2)
	if err != nil {
		return err
	}
	viper.AddConfigPath(callerPath)

	isConfigRead := false
	if IsTesting() {
		// Try config_test.json first
		viper.SetConfigName("config_test")
		if err := viper.ReadInConfig(); err == nil {
			isConfigRead = true
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Logger.Warn("config file config_test.json not found; proceeding with default config")
		} else {
			Logger.Warnf("unable to read config from config_test.json because of error \"%s\"; proceeding with default config", err.Error())
		}
	}
	if !isConfigRead {
		viper.SetConfigName("config")
		if err := viper.ReadInConfig(); err != nil {
			return errors.Newf("unable to read config from config.json").WithCause(err)
		}
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
		} else {
			env = os.Getenv("DEMETER_ENV")
			if env != "" {
				meta.Environment = env
			}
		}
	}

	return nil
}

func GetCallerFilePath(skip int) (string, errors.Error) {
	_, filename, _, ok := runtime.Caller(skip)
	if !ok {
		return "", errors.New("no caller information")
	}
	return path.Dir(filename), nil
}

func GetConfigPath(rootPath string) string {
	return path.Join(rootPath, "config")
}

func GetExecutablePath() (string, errors.Error) {
	exec, err := os.Executable()
	if err != nil {
		return "", errors.Wrap(err)
	}
	return path.Dir(exec), nil
}

func IsTesting() bool {
	return flag.Lookup("test.v") != nil || os.Getenv("ENV") == "test" || os.Getenv("DEMETER_ENV") == "test"
}

func GetMetaConfig(config interface{}) *MetaConfig {
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
			if meta := GetMetaConfig(field.Addr().Interface()); meta != nil {
				return meta
			}
		}
	}
	return nil
}

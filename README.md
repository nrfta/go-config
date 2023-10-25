# go-config ![](https://github.com/nrfta/go-config/workflows/CI/badge.svg)

This package loads config information into a struct. It uses [viper](https://github.com/spf13/viper).

## Installation

```sh
go get github.com/nrfta/go-config/v3
```

## Usage

### Example

Consider this example config content

```json
// config/config.json
{
  "Meta": {
    "environment": "development",
    "service_name": "my-app"
  },
  "port": 6002,
  "bugsnag_api_key": "",
  "segment_write_key": "",
  "log_level": "error"
}
```

In the file you want to load the config in do the following:

```go
// config/config.go

import (
    "embed"

	"github.com/nrfta/go-config/v3"
)

// this MyAppConfig struct is the "custom struct" it has the same attributes that mirror the config json above
type MyAppConfig struct {
	Meta config.MetaConfig

	Port            int    `mapstructure:"port"`
	BugsnagAPIKey   string `mapstructure:"bugsnag_api_key"`
	LogLevel        string `mapstructure:"log_level"`
	PrometheusPort  int    `mapstructure:"prometheus_port"`
	SegmentWriteKey string `mapstructure:"segment_write_key"`
}

//go:embed config.json config_test.json
var fs embed.FS

// the definition of our two paramaters
var (
	Config MyAppConfig
)

func init() {
	err := config.Load(fs, &Config) // now the config data has been loaded into Config
	if err != nil {
		panic(err)
	}
}
```

## License

This project is licensed under the [MIT License](LICENSE.md).

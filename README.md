# go-config ![](https://github.com/neighborly/go-config/workflows/CI/badge.svg)

This package loads config information into a struct. It uses [viper](https://github.com/spf13/viper).

## Installation

```sh
go get github.com/neighborly/go-config
```

## Usage

To load the config data into a struct you will need two parameters

1) a variable of type **box**:   To get the box type you will need to import ``` github.com/gobuffalo/packr``` . The box variable will hold the config data in the binary
2) a variable of type **"customStruct"** where customStruct is a struct you define to mirror the key values of your config

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
	"github.com/gobuffalo/packr"
	"github.com/neighborly/go-config"
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

// the definition of our two paramaters
var (
	Config MyAppConfig
)

func init() {
	err := config.Load(packr.NewBox("."), &Config) // now the config data has been loaded into appConfig
	if err != nil {
		panic(err)
	}
}
```

## License

This project is licensed under the [MIT License](LICENSE.md).

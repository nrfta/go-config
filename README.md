
# go-config

This package loads config information into a struct 


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
// path: folder/config.json 
{
  "Meta": {
    "environment": "test",
    "service_name": "demand-api"
  },
  "port": 6002,
  "bugsnag_api_key": "",
  "segment_write_key": "",
  "log_level": "error"
}
``` 

In the file you want to load the config in do the following: 

```go 
// path folder/file.go 
import (
	"github.com/gobuffalo/packr"
	"github.com/neighborly/go-config"
)
// this DemandAPI struct is the "custom struct" it has the same attributes that mirror the config json above
type DemandAPI struct {
	Meta config.MetaConfig

	Port            int    `mapstructure:"port"`
	BugsnagAPIKey   string `mapstructure:"bugsnag_api_key"`
	LogLevel        string `mapstructure:"log_level"`
	PrometheusPort  int    `mapstructure:"prometheus_port"`
	SegmentWriteKey string `mapstructure:"segment_write_key"`
}

// the definition of our two paramaters
var (
	testBox    = packr.NewBox("path to config data")
	testConfig DemandAPI
)

config.Load(testBox, &testConfig) // now the config data has been loaded into testConfig 
```



## License

This project is licensed under the [MIT License](LICENSE.md).





# config

This package is used for loading the config 

Full unit test implemented 

### Setup
```

import (
    ...
    "github.com/fastretailing/fr-price-common-pkg/config"
    ...
)

const (
	Prod  Env = "prod"
	Stage     = "stage"
	Dev       = "dev"
)

func (e Env) IsValid() bool {
	switch string(e) {
	case "prod", "stage", "dev":
		return true
	}
	return false
}

var configor *config.Config = Default(&envHandler)

type structConfig struct {
	Debug bool   `config:"default:true"`
	Host  string `config:"default:xxxxx"`
}

func (s structConfig) IsValid() bool {
	return true
}

type structConfig2 struct {
	Debug2 bool   `config:"default:true"`
	Host2  string `config:"default:xxxxx"`
}


type testSuccessCase struct {
	TestString  string  `config:"default:test;"`                    // default value as test  convert to string 
	TestInt     int     `config:"default:1;"`                       // default value as 1  convert to int 
	TestInt16   int16   `config:"default:1;"`                       // default value as 1  convert to int16 
	TestInt32   int32   `config:"default:1;"`                       // default value as 1  convert to int32 
	TestInt64   int64   `config:"default:1;"`                       // default value as 1  convert to int64 
	TestUint    uint    `config:"default:1;"`                       // default value as 1  convert to uint 
	TestUint16  uint16  `config:"default:1;"`                       // default value as 1  convert to uint16 
	TestUint32  uint32  `config:"default:1;"`                       // default value as 1  convert to uint32 
	TestUint64  uint64  `config:"default:1;"`                       // default value as 1  convert to uint64 
	TestFloat32 float32 `config:"default:1;"`                       // default value as 1  convert to float32 
	TestFloat64 float64 `config:"default:1;"`                       // default value as 1  convert to float64 
	GoPath      string  `config:"default:1;"`
	GoPathAlias string  `config:"name:GOPATHTEST;default:1;"`       // load key alias to GOPATHTEST
	GoProxy114  string  `config:"name:GOPROXY114TEST;"`             // load key alias to GOPROXY114TEST
	Mandantory  string  `config:"name:mmmmmm;optional:true"`        // optional true , this value can be empty , default is mandantory
	Env         Env     `config:"default:prod;"`                    // support enum
                                                                    // implement config.Field interface with func IsValid()bool 
                                                                    // the config loader will raise err when IsValid check failed 
	DebugFalse  bool    `config:"default:false;"`
	DebugTrue   bool    `config:"default:true;"`
	Sfull       structConfig  `config:"name:sfull;default:{\"name\":\"test\"}`     // support embedded for struct 
	S           structConfig                                        // support embedded for struct 
	SPtr        *structConfig2                                      // support embedded for struct pointer
}

func main(){
    err := config.Default(&envHandler).Load(&c)
    ...
}
	
```


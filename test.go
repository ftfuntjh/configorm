package cofigorm

import "github.com/Unknwon/goconfig"

type Configuration struct {
	ApiConfig   ApiConfiguration
	MysqlConfig MysqlConfiguration
}

type ApiConfiguration struct {
	AccessKey              string `section:"exchange" name:"access-key"`
	SecretKey              string `section:"exchange" name:"secret-key"`
	EnablePrivateSignature bool   `section:"exchange" name:"enable-private-signature" default:"false"`
	PrivateKeyPrime256     string `section:"exchange" name:"private-key-prime256" omit:"true"`
	MarketUrl              string `section:"exchange" name:"market-url"`
	TradeUrl               string `section:"exchange" name:"trade-url"`
	HostName               string `section:"exchange" name:"host-name"`
}

type MysqlConfiguration struct {
	URL string `section:"mysql" name:"url"`
}

const (
	SectionKey     = "section"
	NameKey        = "name"
	DefaultKey     = "default"
	OmitKey        = "omit"
	DefaultSection = goconfig.DEFAULT_SECTION
)

func init() {
	ctg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		panic(err)
	}

	Configure := Configuration{}

	err = Unmarshal(ctg, &Configure)

	if err != nil {
		panic(err)
	}
}

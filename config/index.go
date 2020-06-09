package config

import (
	"os"
)

type Config struct {
	Mode       string
	Connection string
	HMACKey    string
}

func Load() *Config {
	conf := &Config{}
	conf.Mode = os.Getenv("MODE")
	if conf.Mode == "" {
		conf.Mode = "dev"
	}
	conf.Connection = os.Getenv("DB_URL")
	if conf.Connection == "" {
		conf.Connection = "root:password@tcp(poomdv.c52jeww5mzql.ap-southeast-1.rds.amazonaws.com:3306)/hangout?parseTime=true"
	}
	conf.HMACKey = os.Getenv("HMAC_KEY")
	if conf.HMACKey == "" {
		conf.HMACKey = "secret"
	}

	return conf
}

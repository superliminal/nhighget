package config

import (
	"flag"

	"github.com/BurntSushi/toml"
)

type config struct {
	Db         *Db
	IgnoreList []string `toml:"ignore_list"`
	IgnoreIds  []int    `toml:"ignore_ids"`
}

type Db struct {
	//Host     string
	//Port     int
	//Username string
	//Password string
	//Database string
	ConnectionString string `toml:"connection_string"`
}

var File = flag.String("config", "nhighget.conf", "config file path")
var Conf config

func Init(fileLocation string) error {
	_, err := toml.DecodeFile(fileLocation, &Conf)
	return err
}

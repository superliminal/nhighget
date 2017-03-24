package config

import (
	"flag"
	"fmt"

	"github.com/BurntSushi/toml"
)

type config struct {
	Db         *Db
	IgnoreList []string `toml:"ignore_list"`
	IgnoreIds  []int    `toml:"ignore_ids"`
}

type Db struct {
	Username         string
	Password         string
	Host             string
	Port             int
	Database         string
	ConnectionString string `toml:"connection_string"`
}

func (db Db) GetConnectionString() string {
	if len(db.ConnectionString) > 0 {
		return fmt.Sprintf("%v?charset=utf8mb4&parseTime=true", db.ConnectionString)
	}
	return fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Database,
	)
}

var File = flag.String("config", "nhighget.conf", "config file path")
var Conf config

func Init(fileLocation string) error {
	_, err := toml.DecodeFile(fileLocation, &Conf)
	return err
}

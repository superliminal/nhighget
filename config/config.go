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
	username         string
	password         string
	host             string
	port             int
	database         string
	connectionString string `toml:"connection_string"`
}

func (db Db) GetConnectionString() string {
	if len(db.connectionString > 0) {
		return fmt.Sprintf("%v?charset=utf8mb4&parseTime=true", db.connectionString)
	}
	return fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true",
		db.username,
		db.password,
		db.host,
		db.port,
		db.database,
	)
}

var File = flag.String("config", "nhighget.conf", "config file path")
var Conf config

func Init(fileLocation string) error {
	_, err := toml.DecodeFile(fileLocation, &Conf)
	return err
}

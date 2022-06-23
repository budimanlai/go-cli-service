package services

import (
	"log"

	"github.com/eqto/config"
	"github.com/eqto/dbm"
	_ "github.com/eqto/dbm/driver/mysql"
)

var (
	db *dbm.Connection
)

func OpenDb() {
	hostname := config.Get(`database.hostname`)
	username := config.Get(`database.username`)
	password := config.Get(`database.password`)
	port := config.GetInt(`database.port`)
	name := config.Get(`database.database`)

	cn, e := dbm.Connect(dbm.Config{
		DriverName: "mysql",
		Hostname:   hostname,
		Port:       port,
		Username:   username,
		Password:   password,
		Name:       name,
	})
	if e != nil {
		log.Fatal(e)
	}
	db = cn
}

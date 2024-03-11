package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ly1999-hub/simplebank/api"
	db "github.com/ly1999-hub/simplebank/sqlc"
	"github.com/ly1999-hub/simplebank/util"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:Huuly1999@localhost:5432/simple_bank?sslmode=disable"
	address        = "0.0.0.0:8080"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config file env")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to postgresql :", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create NewServer")
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can not start Server")
	}
}

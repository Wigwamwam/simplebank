package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/wigwamwam/simplebank/apic"
	db "github.com/wigwamwam/simplebank/db/sqlc"
	"github.com/wigwamwam/simplebank/util"
)

// chi framework

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	chiStore := db.NewStore(conn)

	apic.NewServer(chiStore)

}

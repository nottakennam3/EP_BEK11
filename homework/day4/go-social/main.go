package main

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"gosocial/store"
	"gosocial/configs"
)

func main() {
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	store, err := store.NewMySQLStorage(cfg)

	if err != nil {
		log.Fatal(err)
	}
	if err = store.Ping(); err != nil {
		log.Fatal(err)
	}
	if err = store.Init(); err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), store)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
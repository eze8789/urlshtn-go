package main

import (
	"flag"
	"fmt"
	"github.com/eze8789/urlshtn-go/pkg/config"
	"github.com/eze8789/urlshtn-go/pkg/database/postgres"
	"github.com/sirupsen/logrus"
)

func main() {
	configFile := flag.String("config", "./configs/app/config.yaml", "Choose configuration file")

	cfg,err := config.GenerateConfig(*configFile)
	if err != nil {
		logrus.Fatalf("could not read configuration file %s: %v", *configFile, err)
	}
	// TODO Remove print statements
	fmt.Println(cfg.Postgres.Port)
	pgConn, err := postgres.NewConn(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password,
		cfg.Postgres.Database, cfg.Postgres.SSLMode)
	if err != nil {
		logrus.Fatalf("could not establish connection: %v", err)
	}
	defer pgConn.DB.Close()

	err = pgConn.CreateTable("./configs/sql/create_url_shortener.sql")
	if err != nil {
		logrus.Fatalf("could not create database: %v", err)
	}

}

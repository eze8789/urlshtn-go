package main

import (
	"flag"
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

	pg, err := postgres.NewConn(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password,

		cfg.Postgres.Database, cfg.Postgres.SSLMode)
	if err != nil {
		logrus.Fatalf("could not establish connection: %v", err)
	}
	defer pg.DB.Close()
}

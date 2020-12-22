package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eze8789/urlshtn-go/pkg/handlers"

	"github.com/eze8789/urlshtn-go/pkg/config"
	"github.com/eze8789/urlshtn-go/pkg/database/postgres"
	"github.com/sirupsen/logrus"
)

const webserverTimeout = 30

func main() {
	logger := logrus.New()

	configFile := flag.String("config", "./configs/app/config.yaml", "Choose configuration file")

	cfg, err := config.GenerateConfig(*configFile)
	if err != nil {
		logger.Fatalf("could not read configuration file %s: %v", *configFile, err)
	}

	pg, err := postgres.NewConn(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password,
		cfg.Postgres.Database, cfg.Postgres.SSLMode)
	if err != nil {
		logger.Fatalf("could not establish connection: %v", err)
	}
	defer pg.Close()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.WebServer.Addr, cfg.WebServer.Port),
		Handler: handlers.Routes(pg),
	}

	//	gracefully shutdown webserver
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		logger.Info("shutting down webserver")
		ctx, cancel := context.WithTimeout(context.Background(), webserverTimeout*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf("could not gracefully shutdown webserver: %v", err)
		}
		close(done)
	}()

	// Start webserver goroutine
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		logger.Infof("starting webserver on: %v", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("could not start webserver on: %v", srv.Addr)
		}
		wg.Done()
	}()
	wg.Wait()

	<-done
	logger.Info("webserver stopped")
}

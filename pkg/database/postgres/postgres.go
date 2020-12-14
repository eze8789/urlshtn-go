package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func NewConn(host, port, username, password, database, ssloption string) (*postgres, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host, port, username, password, database, ssloption)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgres{db}, nil
}

func (d *postgres) CreateTable(filepath string) {}
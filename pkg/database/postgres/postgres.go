package postgres

import (
	"database/sql"
	"fmt"
	"github.com/eze8789/urlshtn-go/pkg/database/models"
	_ "github.com/lib/pq"
	"io/ioutil"
)

type postgres struct {
	DB *sql.DB
}

func NewConn(host, port, username, password, database, ssloption string) (*postgres, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host, port, username, password, database, ssloption)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgres{db}, nil
}

func (d *postgres) CreateTable(filepath string) error {
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = d.DB.Exec(string(f))
	if err != nil {
		return err
	}
	return nil
}

func (d *postgres) Insert(u string) (int, error) {
	stmt := `INSERT INTO url_shortener(url_address) VALUES($1) RETURNING url_id;`
	var id int

	err := d.DB.QueryRow(stmt, u).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// TODO List all
func (d *postgres) List() ([]*models.UrlShort, error) {return nil, nil}

// TODO List all urls
func (d *postgres) ListUrls() ([]string, error) {return nil, nil}

// TODO Retrieve one record
func (d *postgres) Retrieve(u string) (*models.UrlShort, error) {return nil, nil}

//TODO Update visit counts when visited
func (d *postgres) Update(enc string) error {return nil}

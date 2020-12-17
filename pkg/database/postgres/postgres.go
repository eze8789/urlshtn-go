package postgres

import (
	"database/sql"
	"fmt"

	"github.com/eze8789/urlshtn-go/pkg/database/models"

	"io/ioutil"
	"strings"

	// Import to open DB Connection here
	_ "github.com/lib/pq"
)

// Postgres give access to the db struct to execute methods in the handlers
type Postgres struct {
	DB *sql.DB
}

// NewConn establish connection against a postgres DB
func NewConn(host, port, username, password, database, ssloption string) (*Postgres, error) {
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

	err = createTable(db, "./configs/sql/create_url_shortener.sql")
	if err != nil {
		return nil, err
	}
	return &Postgres{db}, nil
}

func createTable(db *sql.DB, filepath string) error {
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(f))
	if err != nil {
		return err
	}
	return nil
}

// Insert a record in the DB
func (d *Postgres) Insert(u string) (int, error) {
	stmt := `INSERT INTO url_shortener(url_address) VALUES($1) RETURNING url_id;`
	if !strings.Contains(u, "https://") {
		u = "https://" + u
	}

	var id int
	err := d.DB.QueryRow(stmt, u).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// List return a slice of UrlShort struct
func (d *Postgres) List() ([]*models.URLShort, error) {
	stmt := `SELECT * FROM url_shortener;`
	rows, err := d.DB.Query(stmt)
	if err == rows.Err() {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*models.URLShort
	for rows.Next() {
		url := &models.URLShort{}
		err = rows.Scan(&url.ID, &url.URLAddress, &url.VisitCounts)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

// Retrieve return a record from the DB
func (d *Postgres) Retrieve(u string) (*models.URLShort, error) {
	stmt := `SELECT * FROM url_shortener WHERE url_address=$1;`
	row := d.DB.QueryRow(stmt, u)
	url := models.URLShort{}

	err := row.Scan(&url.ID, &url.URLAddress, &url.VisitCounts)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
}

// Update change the visit count from a record in the DB
// TODO Update visit counts when visited
func (d *Postgres) Update(enc string) error { return nil }

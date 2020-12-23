package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/eze8789/urlshtn-go/pkg/database"
	"github.com/eze8789/urlshtn-go/pkg/database/models"

	// Import to open DB Connection here
	_ "github.com/lib/pq"
)

// Postgres give access to the db struct to execute methods in the handlers
type postgres struct {
	DB *sql.DB
}

// NewConn establish connection against a postgres DB
func NewConn(host, port, username, password, dbase, ssloption string) (database.Storage, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbase, ssloption)
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
	return &postgres{db}, nil
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
func (p *postgres) Insert(url string) (*int, error) {
	var id int
	stmt := `INSERT INTO url_shortener(url_address) VALUES($1) RETURNING url_id;`
	if !strings.Contains(url, "https://") {
		url = "https://" + url
	}

	err := p.DB.QueryRow(stmt, url).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// List return a slice of UrlShort struct
func (p *postgres) List() (*[]models.URLShort, error) {
	stmt := `SELECT * FROM url_shortener;`
	rows, err := p.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.URLShort
	for rows.Next() {
		url := models.URLShort{}
		err = rows.Scan(&url.ID, &url.URLAddress, &url.VisitCounts)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &urls, nil
}

// Retrieve return a record from the DB
func (p *postgres) RetrieveURL(u string) (string, error) {
	var url string
	stmt := `UPDATE url_shortener SET visit_counts = visit_counts + 1 WHERE url_id= $1 RETURNING url_address;`

	err := p.DB.QueryRow(stmt, u).Scan(&url)
	if err == sql.ErrNoRows {
		return "", err
	}
	if err != nil {
		return "", err
	}

	return url, nil
}

func (p *postgres) RetrieveInfo(u string) (*models.URLShort, error) {
	stmt := `SELECT * FROM url_shortener WHERE url_id = $1;`
	row := p.DB.QueryRow(stmt, u)
	url := models.URLShort{}

	err := row.Scan(&url.ID, &url.URLAddress, &url.VisitCounts)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	fmt.Println(url)
	return &url, nil
}

func (p *postgres) Close() error { return p.DB.Close() }

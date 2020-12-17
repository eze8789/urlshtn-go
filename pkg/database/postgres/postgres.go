package postgres

import (
	"database/sql"
	"fmt"
	"github.com/eze8789/urlshtn-go/pkg/database/models"
	_ "github.com/lib/pq"
	"io/ioutil"
	"strings"
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

func (d *postgres) List() ([]*models.UrlShort, error) {
	stmt := `SELECT * FROM url_shortener;`
	rows, err := d.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*models.UrlShort
	for rows.Next() {
		url := &models.UrlShort{}
		err = rows.Scan(&url.ID, &url.URLAddress, &url.VisitCounts)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (d *postgres) Retrieve(u string) (*models.UrlShort, error) {
	stmt := `SELECT * FROM url_shortener WHERE url_address=$1;`
	row := d.DB.QueryRow(stmt, u)
	url := models.UrlShort{}

	err := row.Scan(&url.ID, &url.URLAddress, &url.VisitCounts)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
}

//TODO Update visit counts when visited
func (d *postgres) Update(enc string) error {return nil}

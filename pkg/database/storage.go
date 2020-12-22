package database

import (
	"github.com/eze8789/urlshtn-go/pkg/database/models"
)

type Storage interface {
	Insert(string) (*int, error)
	List() ([]*models.URLShort, error)
	RetrieveURL(u string) (string, error)
	RetrieveInfo(string) (*models.URLShort, error)
	Update(string) error
	Close() error
}

package storage

import (
	"time"
)

type Storage interface {
	Save(p *Page) error
}

type Page struct {
	Text       string
	UserName   string
	CreatedAt  time.Time
	PictureUrl string
}

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Clients ClientModel
	Files   FileModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Clients: ClientModel{DB: db},
		Files:   FileModel{DB: db},
	}
}

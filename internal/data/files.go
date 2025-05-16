package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type File struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	Category  string    `json:"category"`
	ClientID  int       `json:"client_id"`
}

type FileModel struct {
	DB *sql.DB
}

func (m FileModel) Insert(files []File) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	query := `
        INSERT INTO files (file_name, file_path, category, client_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range files {
		file := &files[i]
		args := []any{file.FileName, file.FilePath, file.Category, file.ClientID}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := stmt.QueryRowContext(ctx, args...).Scan(&file.ID, &file.CreatedAt)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				return fmt.Errorf("query failed: %v; rollback failed: %v", err, rbErr)
			}
			return fmt.Errorf("query failed: %w", err)
		}
	}

	return nil
}

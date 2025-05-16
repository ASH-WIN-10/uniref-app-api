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

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("commit failed: %v; rollback failed: %v", err, rbErr)
		}
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (m FileModel) Get(clientID int) ([]File, error) {
	query := `
        SELECT id, created_at, file_name, file_path, category, client_id
        FROM files
        WHERE client_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, clientID)
	if err != nil {
		return nil, err
	}

	files := []File{}
	for rows.Next() {
		var file File

		err := rows.Scan(
			&file.ID,
			&file.CreatedAt,
			&file.FileName,
			&file.FilePath,
			&file.Category,
			&file.ClientID,
		)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

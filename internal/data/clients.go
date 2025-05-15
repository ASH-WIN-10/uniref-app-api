package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Client struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"-"`
	CompanyName string    `json:"company_name"`
	ClientName  string    `json:"client_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
}

type ClientModel struct {
	DB *sql.DB
}

func (m *ClientModel) Insert(client *Client) error {
	query := `
        INSERT INTO clients (company_name, client_name, email, phone)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`

	args := []any{client.CompanyName, client.ClientName, client.Email, client.Phone}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&client.ID, &client.CreatedAt)
}

func (m *ClientModel) Get(id int) (*Client, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, company_name, client_name, email, phone
        FROM clients
        WHERE id = $1`

	var client Client

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&client.ID,
		&client.CreatedAt,
		&client.CompanyName,
		&client.ClientName,
		&client.Email,
		&client.Phone,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &client, nil
}

func (m *ClientModel) Update(client *Client) error {
	return nil
}

func (m *ClientModel) Delete(id int) error {
	return nil
}

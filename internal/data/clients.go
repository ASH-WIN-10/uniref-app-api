package data

import "time"

type Client struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"-"`
	CompanyName string    `json:"company_name"`
	ClientName  string    `json:"client_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
}

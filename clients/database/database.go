package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Credentials struct {
	Username string
	Password string
	DBName   string
	Host     string
	SSLMode  *string
}

func (d *Credentials) SourceString() string {
	source := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s",
		d.Username, d.Password, d.DBName, d.Host,
	)

	if d.SSLMode != nil {
		source += fmt.Sprintf(" sslmode=%s", *d.SSLMode)
	}

	return source
}

type Client struct {
	*sqlx.DB
}

func (d *Client) Check(ctx context.Context) error {
	return d.PingContext(ctx)
}

func New(ctx context.Context, driverName string, creds *Credentials) (*Client, error) {
	connection, err := sqlx.Open(driverName, creds.SourceString())
	if err != nil {
		return nil, err
	}

	if err := connection.Ping(); err != nil {
		return nil, err
	}

	return &Client{
		connection,
	}, nil
}

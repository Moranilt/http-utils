package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Credentials struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (d *Credentials) SourceString(production bool) string {
	if !production {
		return fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=disable",
			d.Username, d.Password, d.DBName, d.Host,
		)
	}
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s",
		d.Username, d.Password, d.DBName, d.Host,
	)
}

type Client struct {
	*sqlx.DB
}

func (d *Client) Check(ctx context.Context) error {
	return d.PingContext(ctx)
}

func New(ctx context.Context, driverName string, creds *Credentials, production bool) (*Client, error) {
	connection, err := sqlx.Open(driverName, creds.SourceString(production))
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

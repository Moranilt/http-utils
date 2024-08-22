package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makePointer[T any](val T) *T {
	return &val
}

func TestCredentials_SourceString(t *testing.T) {
	t.Run("with sslmode", func(t *testing.T) {
		creds := Credentials{
			Username: "test",
			Password: "<PASSWORD>",
			DBName:   "test",
			Host:     "test",
			SSLMode:  makePointer("test"),
		}

		expected := "user=test password=<PASSWORD> dbname=test host=test sslmode=test"
		assert.Equal(t, expected, creds.SourceString())
	})

	t.Run("without sslmode", func(t *testing.T) {
		creds := Credentials{
			Username: "test",
			Password: "<PASSWORD>",
			DBName:   "test",
			Host:     "test",
			SSLMode:  nil,
		}
		expected := "user=test password=<PASSWORD> dbname=test host=test"
		assert.Equal(t, expected, creds.SourceString())
	})
}

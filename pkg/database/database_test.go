package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConnectionString(t *testing.T) {
	tests := []struct {
		dbConfig Config
		want     string
	}{
		{
			Config{
				Host:     "amountly-not-my-localhost",
				Port:     1,
				Username: "someone",
				Password: "1234",
				Name:     "not-my-sql",
			},
			"postgres://someone:1234@amountly-not-my-localhost:1/not-my-sql?sslmode=disable",
		},
	}

	for _, test := range tests {
		result := buildConnectionString(test.dbConfig)
		assert.Equal(t, result, test.want)
	}
}

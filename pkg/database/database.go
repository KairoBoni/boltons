package database

import (
	"fmt"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //database driver
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Config sets the configuration for the database
type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

func genSchema() string {
	dbTimezone := `SET TIMEZONE TO 'America/Sao_Paulo';`

	nfeAmountTable := `
		CREATE TABLE IF NOT EXISTS nfe_amount (
		access_key TEXT PRIMARY KEY,
		amount TEXT NOT NULL
	);`

	// We only do this action for demonstration purposes
	// if you do not delete what is in the database
	// when it runs again we will have problems with duplicate primary keys
	deletePersistentData := `
		DELETE FROM nfe_amount;
	`

	return dbTimezone + nfeAmountTable + deletePersistentData
}

func buildConnectionString(cfg Config) string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}

// CreateDB creates a new connection pool with the database and applies the database schema.
// The Postgres configurations can be found at docker-compose.yml in postgres container
func CreateDB(cfg Config) (*sqlx.DB, error) {
	const maxAttempts = 6
	var (
		db  *sqlx.DB
		err error
	)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		db, err = sqlx.Open("postgres", buildConnectionString(cfg))
		if err == nil {
			break
		} else if attempt == maxAttempts-1 {
			log.Info().Msgf("Failed to connect to database, max attempts reached: %v", err)
			return nil, errors.Wrap(err, "failed to connect to database")
		}

		wait := time.Duration(math.Exp2(float64(attempt))) * time.Second
		log.Info().Msgf("Failed to connect to database. Retrying in %v...", wait)
		time.Sleep(wait)
	}

	_, err = db.Exec(genSchema())
	if err != nil {
		return nil, fmt.Errorf("failed to exec database schema: %+v", err)
	}
	log.Info().Msg("connected to database")
	return db, nil
}

package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	dbConfigFilepath := os.Getenv("CONFIG_DB_FILEPATH")
	if dbConfigFilepath == "" {
		log.Error().Msgf("Missing environment variable CONFIG_DB_FILEPATH")
	}

	s, err := NewServer(dbConfigFilepath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create a new store from database")
	}

	if err := s.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start api")
	}

}

package database

import (
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

type Store struct {
	db *sqlx.DB
}

func getConfig(filepath string) (*Config, error) {
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewStore(cfgFilepath string) (*Store, error) {
	cfg, err := getConfig(cfgFilepath)
	if err != nil {
		return nil, err
	}

	db, err := CreateDB(*cfg)
	if err != nil {
		return nil, err
	}
	return &Store{db}, nil
}

func (s *Store) Ping() error {
	return s.db.Ping()
}

func (s *Store) InsertNfeTotal(accessKey, total string) error {
	_, err := s.db.Exec(insertNFETotal, accessKey, total)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetNfeTotal(accessKey string) (string, error) {
	var total string

	row := s.db.QueryRow(getNFETotal, accessKey)
	if err := row.Scan(&total); err != nil {
		return "", err
	}

	return total, nil
}

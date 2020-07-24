package database

import (
	"database/sql"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

//Store implements the StoreInterface
type Store struct {
	db *sqlx.DB
}

//StoreInterface interface to all actions in database
type StoreInterface interface {
	InsertNfeAmount(accessKey, amount string) error
	GetNfeAmount(accessKey string) (string, error)
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

//NewStore create a new store with the seted config
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

//Ping just test the connection with database
func (s *Store) Ping() error {
	return s.db.Ping()
}

//InsertNfeAmount insert a new line with the access key and the amount of nfe value
func (s *Store) InsertNfeAmount(accessKey, amount string) error {
	_, err := s.db.Exec(insertNFEAmount, accessKey, amount)
	if err != nil {
		return err
	}
	return nil
}

//GetNfeAmount get the amount nfe value associated with the set access key
func (s *Store) GetNfeAmount(accessKey string) (string, error) {
	var amount string

	row := s.db.QueryRow(getNFEAmount, accessKey)
	if err := row.Scan(&amount); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return amount, nil
}

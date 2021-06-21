package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	config *Config
	db *sql.DB
	advRepository *AdvRepository
	refrepository *refrepository
}

func NewStore(config *Config) *Store{
	s := Store{
		config: config,
	}
	return &s
}

func (s *Store) Open() error{
	db, err := sql.Open("postgres", s.config.DataBaseUrl)
	if (err != nil) {
		return err
	}

	if err := db.Ping(); err != nil{
		return err
	}
	s.db = db
	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Adv() *AdvRepository{
	if s.advRepository != nil {
		return s.advRepository
	}

	s.advRepository = &AdvRepository{
		store:s,
	}
	return s.advRepository
}

func (s *Store) Ref() *refrepository{
	if s.refrepository != nil {
		return s.refrepository
	}

	s.refrepository = &refrepository{
		store:s,
	}
	return s.refrepository
}

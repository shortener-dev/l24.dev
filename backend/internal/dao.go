package internal

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	InsertShortQuery = "INSERT INTO urls (redirect_path, scheme, host, path, query) VALUES ($1,$2,$3,$4,$5)"
	GetShortQuery    = "SELECT redirect_path, scheme, host, path, query FROM urls WHERE redirect_path=$1"
)

type ShortDAO interface {
	InsertShort(short Short) error
	GetShort(redirect_path string) (*Short, error)
}

func NewShortPostgresDao(db *sql.DB, driver string) *ShortPostgresDAO {
	return &ShortPostgresDAO{db: db, driver: driver}
}

type ShortPostgresDAO struct {
	db     *sql.DB
	driver string
}

func (s *ShortPostgresDAO) InsertShort(short Short) error {
	db := sqlx.NewDb(s.db, s.driver)

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		InsertShortQuery,
		short.RedirectPath,
		short.Scheme,
		short.Host,
		short.Path,
		short.Query,
	)
	if err != nil {
		log.Printf("error inserting url: %v", err)
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Printf("error rolling back transaction: %v ", rollbackErr)
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("error committing transaction: %q", err)
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Printf("error rolling back transaction: %v ", rollbackErr)
			return rollbackErr
		}
		return err
	}

	return nil
}

func (s *ShortPostgresDAO) GetShort(redirect_path string) (*Short, error) {
	db := sqlx.NewDb(s.db, s.driver)

	row := db.QueryRowx(GetShortQuery, redirect_path)

	var short Short
	err := row.StructScan(&short)
	if err != nil {
		return nil, err
	}

	return &short, nil
}

package shortener

import (
	"context"
	"database/sql"

	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres Driver
)

const (
	InsertShortQuery          = "INSERT INTO urls (redirect_path, scheme, host) VALUES ($1,$2,$3)"
	InsertShortWithPathQuery  = "INSERT INTO urls (redirect_path, scheme, host, path) VALUES ($1,$2,$3,$4)"
	InsertShortWithQueryQuery = "INSERT INTO urls (redirect_path, scheme, host, query) VALUES ($1,$2,$3,$4)"
	InsertShortWithAllQuery   = "INSERT INTO urls (redirect_path, scheme, host, path, query) VALUES ($1,$2,$3,$4,$5)"
	GetShortQuery             = "SELECT redirect_path, scheme, host, path, query FROM urls WHERE redirect_path=$1"
)

type ShortDAO interface {
	InsertShort(ctx context.Context, short Short) error
	GetShort(ctx context.Context, redirect_path string) (*Short, error)
}

func NewShortPostgresDao(db *sql.DB, driver string) *ShortPostgresDAO {
	return &ShortPostgresDAO{db: db, driver: driver}
}

type ShortPostgresDAO struct {
	db     *sql.DB
	driver string
}

func (s *ShortPostgresDAO) InsertShort(ctx context.Context, short Short) error {
	db := sqlx.NewDb(s.db, s.driver)

	var err error

	switch {
	case isNilOrEmptyString(short.Path) && isNilOrEmptyString(short.Query):
		err = executeTransaction(
			ctx,
			*db,
			InsertShortQuery,
			short.RedirectPath,
			short.Scheme,
			short.Host,
		)
	case isNilOrEmptyString(short.Path):
		err = executeTransaction(
			ctx,
			*db,
			InsertShortWithQueryQuery,
			short.RedirectPath,
			short.Scheme,
			short.Host,
			short.Query,
		)
	case isNilOrEmptyString(short.Query):
		err = executeTransaction(
			ctx,
			*db,
			InsertShortWithPathQuery,
			short.RedirectPath,
			short.Scheme,
			short.Host,
			short.Path,
		)
	default:
		err = executeTransaction(
			ctx,
			*db,
			InsertShortWithAllQuery,
			short.RedirectPath,
			short.Scheme,
			short.Host,
			short.Path,
			short.Query,
		)
	}

	return err
}

func (s *ShortPostgresDAO) GetShort(ctx context.Context, redirect_path string) (*Short, error) {
	db := sqlx.NewDb(s.db, s.driver)

	row := db.QueryRowxContext(ctx, GetShortQuery, redirect_path)

	var short Short
	err := row.StructScan(&short)
	if err != nil {
		return nil, err
	}

	return &short, nil
}

func isNilOrEmptyString(text *string) bool {
	if text == nil {
		return true
	}

	if *text == "" {
		return true
	}

	return false
}

func executeTransaction(ctx context.Context, db sqlx.DB, query string, args ...interface{}) error {
	var err error
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = multierror.Append(err, rollbackErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = multierror.Append(err, rollbackErr)
		}
		return err
	}

	return nil
}

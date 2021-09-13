package shortener

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres Driver
)

const (
	InsertShortQuery = "INSERT INTO urls (%v) VALUES (%v)"
	GetShortQuery    = "SELECT redirect_path, scheme, host, path, query, fragment FROM urls WHERE redirect_path=$1"
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

	query := s.buildInsertQuery(short)
	return executeTransaction(ctx, *db, query)
}

func (s *ShortPostgresDAO) GetShort(ctx context.Context, redirect_path string) (*Short, error) {
	db := sqlx.NewDb(s.db, s.driver)

	var short Short
	err := db.GetContext(ctx, &short, GetShortQuery, redirect_path)
	if err != nil {
		return nil, err
	}

	return &short, nil
}

func (s *ShortPostgresDAO) buildInsertQuery(short Short) string {
	columns := []string{"redirect_path", "scheme", "host"}

	values := []string{
		addSingleQuotes(short.RedirectPath),
		addSingleQuotes(short.Scheme),
		addSingleQuotes(short.Host),
	}

	if !isNilOrEmptyString(short.Path) {
		columns = append(columns, "path")
		values = append(values, addSingleQuotes(*short.Path))
	}

	if !isNilOrEmptyString(short.Query) {
		columns = append(columns, "query")
		values = append(values, addSingleQuotes(*short.Query))
	}

	if !isNilOrEmptyString(short.Fragment) {
		columns = append(columns, "fragment")
		values = append(values, addSingleQuotes(*short.Fragment))
	}

	columnsString := strings.Join(columns, ", ")
	valuesString := strings.Join(values, ", ")

	return fmt.Sprintf(InsertShortQuery, columnsString, valuesString)
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

func addSingleQuotes(s string) string {
	return "'" + s + "'"
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

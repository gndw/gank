package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/db"
	"github.com/gndw/gank/services/secret"
	"github.com/gndw/gank/services/utils/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Service struct {
	db  *sqlx.DB
	log log.Service
}

func NewSqlx(lc model.Lifecycle, secret secret.Service, log log.Service) (db.Service, error) {
	ins := &Service{
		log: log,
	}

	sdb := secret.Database
	host, port, user, pw, dbn := sdb.Host, sdb.Port, sdb.User, sdb.Password, sdb.DBName
	dataSource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s ", host, port, user, pw, dbn)

	if !functions.IsAllAboveZero(port) || !functions.IsAllNonEmpty(host, user, pw, dbn) {
		return nil, errors.New("database credentials cannot be empty")
	}

	err := functions.LoggingProcessTime(log, "connect to database", func() error {
		db, err := sqlx.Connect("postgres", dataSource)
		if err != nil {
			return err
		}

		ins.db = db
		return nil
	},
		// Will produce warning log if the process took more than 4 seconds
		functions.WithLoggingProcessTimeLimit(4),
	)
	if err != nil {
		return nil, err
	}

	// Close db connection when error at startup
	lc.AppendOnError(func(ctx context.Context) error {
		return functions.LoggingProcessTime(log, "close database connection", func() error {
			return ins.db.Close()
		})
	})

	// Close db connection when application is shutdown
	lc.Append(model.NewHook(model.Hook{
		OnStop: func(ctx context.Context) error {
			return functions.LoggingProcessTime(log, "close database connection", func() error {
				return ins.db.Close()
			})
		},
	}))

	return ins, nil
}

type Transaction struct {
	tx *sql.Tx
}

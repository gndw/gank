package impl

import (
	"context"
	"database/sql"

	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/services/db"
	"github.com/jmoiron/sqlx"
)

func (s *Service) Ping(ctx context.Context) (err error) {
	return s.db.PingContext(ctx)
}

func (s *Service) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	return s.db.SelectContext(ctx, dest, query, args...)
}

func (s *Service) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	err = s.db.GetContext(ctx, dest, query, args...)
	if err == sql.ErrNoRows {
		err = errorsg.WithOptions(err, errorsg.WithPrivateIdentifier(db.ERROR_IDENTIFIER_DB_GET_NOT_FOUND))
	}
	return err
}

func (s *Service) Exec(ctx context.Context, query string, args ...interface{}) (err error) {
	_, err = s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Service) WithTransaction(ctx context.Context, transaction func(context.Context, db.Transaction) error) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		err = errorsg.InternalErrorWithOptions(err, errorsg.WithPrivateIdentifier(db.ERROR_IDENTIFIER_TX_FAILED))
		return err
	}
	myTx := Transaction{tx: tx}
	err = transaction(ctx, &myTx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			s.log.Errorf("transaction failed and database failed to rollback")
			err = errorsg.InternalErrorWithOptions(err, errorsg.WithPrivateIdentifier(db.ERROR_IDENTIFIER_ROLLBACK_FAILED))
		} else {
			s.log.Warningf("transaction failed and database rollback-ed")
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		err = errorsg.InternalErrorWithOptions(err, errorsg.WithPrivateIdentifier(db.ERROR_IDENTIFIER_TX_FAILED))
		return err
	}

	return nil
}

func (t *Transaction) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	return sqlx.StructScan(rows, dest)
}

func (t *Transaction) Get(ctx context.Context, dests []interface{}, query string, args ...interface{}) (err error) {
	err = t.tx.QueryRowContext(ctx, query, args...).Scan(dests...)
	if err == sql.ErrNoRows {
		err = errorsg.WithOptions(err, errorsg.WithPrivateIdentifier(db.ERROR_IDENTIFIER_DB_GET_NOT_FOUND))
	}
	return err
}

func (t *Transaction) Exec(ctx context.Context, query string, args ...interface{}) (err error) {
	_, err = t.tx.ExecContext(ctx, query, args)
	return err
}

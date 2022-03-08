package db

import "context"

var ERROR_IDENTIFIER_ROLLBACK_FAILED string = "ERROR_IDENTIFIER_ROLLBACK_FAILED"
var ERROR_IDENTIFIER_TX_FAILED string = "ERROR_IDENTIFIER_TX_FAILED"
var ERROR_IDENTIFIER_DB_GET_NOT_FOUND string = "ERROR_IDENTIFIER_DB_GET_NOT_FOUND"

type Service interface {
	Ping(ctx context.Context) (err error)
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)
	Exec(ctx context.Context, query string, args ...interface{}) (err error)
	WithTransaction(ctx context.Context, transaction func(context.Context, Transaction) error) (err error)
	ArgArray(a interface{}) interface{}
}

type Transaction interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)
	Get(ctx context.Context, dests []interface{}, query string, args ...interface{}) (err error)
	Exec(ctx context.Context, query string, args ...interface{}) (err error)
}

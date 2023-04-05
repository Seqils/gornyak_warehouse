package database

import (
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	pgdial "github.com/uptrace/bun/dialect/pgdialect"
	pgd "github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewDBConnection() *bun.DB {
	pgConf := pgd.NewConnector(
		pgd.WithNetwork("tcp"),
		pgd.WithAddr("localhost:5432"),
		pgd.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgd.WithUser("postgres"),
		pgd.WithPassword("postgres"),
		pgd.WithDatabase("gornyak"),
		pgd.WithApplicationName("gornyakWarehouse"),
		pgd.WithDialTimeout(5 *time.Second),
		pgd.WithTimeout(5 * time.Second),
		pgd.WithReadTimeout(5 * time.Second),
		pgd.WithWriteTimeout(5 * time.Second),
	)

	sqlDb := sql.OpenDB(pgConf)

	db := bun.NewDB(sqlDb, pgdial.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}
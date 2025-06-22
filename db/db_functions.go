package db

import (
	"context"
	"log"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	db   *sqlx.DB
	once sync.Once
)

func GetDB(ctx context.Context) *sqlx.DB {
	connStr := "postgres://myuser:mypassword@localhost:5432/mydatabase"

	once.Do(func() {
		var err error
		// Open the database connection
		db, err = sqlx.Open("pgx", connStr)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		// Verify the connection
		if err := db.PingContext(ctx); err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}

		// Set connection pool parameters
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		// Optionally set connection lifetime
		// db.SetConnMaxLifetime(time.Hour)
	})

	return db
}

type Order struct {
	OrderID     int64   `db:"order_id"`
	CustomerID  int64   `db:"customer_id"`
	TotalAmount float64 `db:"total_amount"`
	Status      string  `db:"status"`
}

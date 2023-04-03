package dbconnection

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgPoolConnection *pgxpool.Pool

func newPoolConnections() {
	logger := log.GetLogger()
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	databaseURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		logger.Fatal("Unable to connect to database: %v\n", err)
	}

	pgPoolConnection = dbpool
}

func GetPoolConnections() *pgxpool.Pool {
	var once sync.Once
	once.Do(newPoolConnections)
	return pgPoolConnection
}

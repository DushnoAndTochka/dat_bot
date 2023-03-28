package dbconnection

import (
	"context"
	"os"
	"sync"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgPoolConnection *pgxpool.Pool

func newPoolConnections() {
	logger := log.GetLogger()
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	pgPoolConnection = dbpool
}

func GetPoolConnections() *pgxpool.Pool {
	var once sync.Once
	once.Do(newPoolConnections)
	return pgPoolConnection
}

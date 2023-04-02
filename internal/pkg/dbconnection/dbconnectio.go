package dbconnection

import (
	"context"
	"fmt"
	"sync"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgPoolConnection *pgxpool.Pool

func newPoolConnections() {
	logger := log.GetLogger()
	databaseURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:5432/%s",
		"project",
		"project",
		"host.docker.internal",
		"project",
	)

	logger.Error(databaseURL)

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

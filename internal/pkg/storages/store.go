package storages

import (
	"sync"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/dbconnection"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var storage *Store

type Store struct {
	conn   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func GetStorage() *Store {
	var once sync.Once
	once.Do(func() {
		logger := log.GetLogger()
		conn := dbconnection.GetPoolConnections()
		storage = &Store{
			logger: logger,
			conn:   conn,
		}
	})

	return storage
}

package storages

import (
	"context"
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
	ctx    context.Context
}

func NewStorage(ctx context.Context) *Store {
	var once sync.Once
	once.Do(func() {
		logger := log.GetLogger()
		conn := dbconnection.GetPoolConnections()
		storage = &Store{
			logger: logger,
			conn:   conn,
			ctx:    ctx,
		}
	})

	return storage
}

func GetStorage() *Store {
	return storage
}

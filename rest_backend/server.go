package rest_backend

import (
	"github.com/makeroo/taxi_scout/storage"
	"go.uber.org/zap"
)

type RestServer struct {
	Dao    storage.Datastore
	Logger *zap.SugaredLogger
	Configuration Configuration
}

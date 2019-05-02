package rest

import (
	"github.com/makeroo/taxi_scout/storage"
	"go.uber.org/zap"
)

// Server implements Taxi Scout REST API.
type Server struct {
	Dao           storage.Datastore
	Logger        *zap.SugaredLogger
	Configuration Configuration
}

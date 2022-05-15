package status

import (
	"app/external/database"
	"fmt"
	"github.com/hellofresh/health-go/v4"
	psql "github.com/hellofresh/health-go/v4/checks/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type API struct {
	logger *zerolog.Logger
	db     *sqlx.DB
	config database.Config
}

func NewAPI(logger *zerolog.Logger, db *sqlx.DB, config database.Config) *API {
	return &API{
		logger: logger,
		db:     db,
		config: config,
	}
}

func (api *API) RegisterHandlers(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/status", api.registerHealthChecks())
}

func (api *API) registerHealthChecks() http.HandlerFunc {
	healthStatus, err := health.New()
	if err != nil {
		api.logger.Fatal().Err(err).Msg("Failed to create health container")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		api.config.Host, api.config.Port, api.config.Username, api.config.Password, api.config.Database,
	)
	err = healthStatus.Register(health.Config{
		Name:      "postgresql",
		Timeout:   time.Second * 30,
		SkipOnErr: false,
		Check: psql.New(psql.Config{
			DSN: psqlInfo,
		}),
	})
	if err != nil {
		api.logger.Fatal().Err(err).Msg("Failed to create postgres health check")
	}

	return healthStatus.HandlerFunc
}

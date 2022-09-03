package status

import (
	"fmt"
	"github.com/hellofresh/health-go/v4"
	psql "github.com/hellofresh/health-go/v4/checks/postgres"
	"net/http"
	"time"
)

func (api *API) registerHealthChecks() http.HandlerFunc {
	healthStatus, err := health.New()
	if err != nil {
		api.logger.Log().
			Fatal().
			Err(err).
			Msg("Failed to create health container")
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
		api.logger.Log().
			Fatal().
			Err(err).
			Msg("Failed to create postgres health check")
	}

	return healthStatus.HandlerFunc
}

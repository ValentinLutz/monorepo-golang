package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type StatsGetter interface {
	Stats() sql.DBStats
}

// DatabaseOpts is a reduced variation of the prometheus.Opts
type DatabaseOpts struct {
	Namespace   string
	Subsystem   string
	ConstLabels prometheus.Labels
}

// DatabaseStats contains database statistics from sql.DBStats.
type DatabaseStats struct {
	StatsGetter

	maxOpenConnections *prometheus.Desc // Maximum number of open connections to the database.

	// Pool Status
	openConnections *prometheus.Desc // The number of established connections both in use and idle.
	inUse           *prometheus.Desc // The number of connections currently in use.
	idle            *prometheus.Desc // The number of idle connections.

	// Counters
	waitCount         *prometheus.Desc // The total number of connections waited for.
	waitDuration      *prometheus.Desc // The total time blocked waiting for a new connection.
	maxIdleClosed     *prometheus.Desc // The total number of connections closed due to SetMaxIdleConns.
	maxIdleTimeClosed *prometheus.Desc // The total number of connections closed due to SetConnMaxIdleTime.
	maxLifetimeClosed *prometheus.Desc // The total number of connections closed due to SetConnMaxLifetime.
}

func NewDatabaseStats(statsGetter StatsGetter, opts DatabaseOpts) *DatabaseStats {
	return &DatabaseStats{
		StatsGetter: statsGetter,
		maxOpenConnections: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_max_open_connections"),
			"Maximum number of open connections to the database.",
			nil,
			opts.ConstLabels,
		),
		openConnections: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_open_connections"),
			"The number of established connections both in use and idle.",
			nil,
			opts.ConstLabels,
		),
		inUse: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_in_use"),
			"The number of connections currently in use.",
			nil,
			opts.ConstLabels,
		),
		idle: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_idle"),
			"The number of idle connections.",
			nil,
			opts.ConstLabels,
		),
		waitCount: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_wait_count"),
			"The total number of connections waited for.",
			nil,
			opts.ConstLabels,
		),
		waitDuration: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_wait_duration_seconds"),
			"The total time blocked waiting for a new connection.",
			nil,
			opts.ConstLabels,
		),
		maxIdleClosed: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_max_idle_closed"),
			"The total number of connections closed due to SetMaxIdleConns.",
			nil,
			opts.ConstLabels,
		),
		maxLifetimeClosed: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_max_lifetime_closed"),
			"The total number of connections closed due to SetConnMaxLifetime.",
			nil,
			opts.ConstLabels,
		),
		maxIdleTimeClosed: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "database_max_idletime_closed"),
			"The total number of connections closed due to SetConnMaxIdleTime.",
			nil,
			opts.ConstLabels,
		),
	}

}

// Describe implements the prometheus.Collector interface.
func (databaseStats DatabaseStats) Describe(ch chan<- *prometheus.Desc) {
	ch <- databaseStats.maxOpenConnections
	ch <- databaseStats.openConnections
	ch <- databaseStats.inUse
	ch <- databaseStats.idle
	ch <- databaseStats.waitCount
	ch <- databaseStats.waitDuration
	ch <- databaseStats.maxIdleClosed
	ch <- databaseStats.maxLifetimeClosed
	ch <- databaseStats.maxIdleTimeClosed
}

// Collect implements the prometheus.Collector interface.
func (databaseStats DatabaseStats) Collect(ch chan<- prometheus.Metric) {
	stats := databaseStats.Stats()

	ch <- prometheus.MustNewConstMetric(
		databaseStats.maxOpenConnections,
		prometheus.GaugeValue,
		float64(stats.MaxOpenConnections),
	)
	ch <- prometheus.MustNewConstMetric(
		databaseStats.openConnections,
		prometheus.GaugeValue,
		float64(stats.OpenConnections),
	)
	ch <- prometheus.MustNewConstMetric(
		databaseStats.inUse,
		prometheus.GaugeValue,
		float64(stats.InUse),
	)
	ch <- prometheus.MustNewConstMetric(
		databaseStats.idle,
		prometheus.GaugeValue,
		float64(stats.Idle),
	)
	ch <- prometheus.MustNewConstMetric(
		databaseStats.waitCount,
		prometheus.CounterValue,
		float64(stats.WaitCount),
	)
	ch <- prometheus.MustNewConstMetric(
		databaseStats.waitDuration,
		prometheus.CounterValue,
		float64(stats.WaitDuration.Seconds()),
	)
	ch <- prometheus.MustNewConstMetric(
		databaseStats.maxIdleClosed,
		prometheus.CounterValue,
		float64(stats.MaxIdleClosed),
	)
	ch <- prometheus.MustNewConstMetric(
		databaseStats.maxLifetimeClosed,
		prometheus.CounterValue,
		float64(stats.MaxLifetimeClosed),
	)
	ch <- prometheus.MustNewConstMetric(
		databaseStats.maxIdleTimeClosed,
		prometheus.CounterValue,
		float64(stats.MaxIdleTimeClosed),
	)
}

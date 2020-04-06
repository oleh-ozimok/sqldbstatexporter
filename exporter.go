package sqldbstatxporter

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const subsystem = "conn_pool"

type Statter interface {
	Stats() sql.DBStats
}

type Exporter struct {
	statter                  Statter
	maxOpenConns             *prometheus.Desc
	openConns                *prometheus.Desc
	inUse                    *prometheus.Desc
	idle                     *prometheus.Desc
	waitCount                *prometheus.Desc
	waitDurationMicroseconds *prometheus.Desc
	maxIdleClosed            *prometheus.Desc
	maxLifetimeClosed        *prometheus.Desc
}

func New(statter Statter, namespace string, constantLabels prometheus.Labels) *Exporter {
	return &Exporter{
		statter: statter,
		maxOpenConns: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "max_open_conns"),
			"Maximum number of open connections to the database",
			nil,
			constantLabels,
		),
		openConns: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "open_conns"),
			"The number of established connections both in use and idle",
			nil,
			constantLabels,
		),
		inUse: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "in_use"),
			"The number of connections currently in use",
			nil,
			constantLabels,
		),
		idle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "idle"),
			"The number of idle connections",
			nil,
			constantLabels,
		),
		waitCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wait_count"),
			"The total number of connections waited for",
			nil,
			constantLabels,
		),
		waitDurationMicroseconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wait_duration_microseconds"),
			"The total time blocked waiting for a new connection",
			nil,
			constantLabels,
		),
		maxIdleClosed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "max_idle_closed"),
			"The total number of connections closed due to SetMaxIdleConns",
			nil,
			constantLabels,
		),
		maxLifetimeClosed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "max_lifetime_closed"),
			"The total number of connections closed due to SetConnMaxLifetime",
			nil,
			constantLabels,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.maxOpenConns
	ch <- e.openConns
	ch <- e.inUse
	ch <- e.idle
	ch <- e.waitCount
	ch <- e.waitDurationMicroseconds
	ch <- e.maxIdleClosed
	ch <- e.maxLifetimeClosed
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	s := e.statter.Stats()

	ch <- prometheus.MustNewConstMetric(e.maxOpenConns, prometheus.GaugeValue, float64(s.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(e.openConns, prometheus.GaugeValue, float64(s.OpenConnections))
	ch <- prometheus.MustNewConstMetric(e.inUse, prometheus.GaugeValue, float64(s.InUse))
	ch <- prometheus.MustNewConstMetric(e.idle, prometheus.GaugeValue, float64(s.Idle))
	ch <- prometheus.MustNewConstMetric(e.waitCount, prometheus.GaugeValue, float64(s.WaitCount))
	ch <- prometheus.MustNewConstMetric(e.waitDurationMicroseconds, prometheus.GaugeValue, float64(s.WaitDuration.Microseconds()))
	ch <- prometheus.MustNewConstMetric(e.maxIdleClosed, prometheus.GaugeValue, float64(s.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(e.maxLifetimeClosed, prometheus.GaugeValue, float64(s.MaxLifetimeClosed))
}

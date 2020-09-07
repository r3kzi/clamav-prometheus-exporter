package clamav

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/cfg"
)

type Collector struct {
	config      cfg.Config
	status      *prometheus.Desc
	threadsLive *prometheus.Desc
}

func NewCollector(config cfg.Config) *Collector {
	return &Collector{
		config:      config,
		status:      prometheus.NewDesc("clamav_status", "Shows UP Status", nil, nil),
		threadsLive: prometheus.NewDesc("clamav_threads_live", "Shows live threads", nil, nil),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.status
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	address := fmt.Sprintf("%s:%d", collector.config.ClamAVAddress, collector.config.ClamAVPort)
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.CounterValue, ping(address))
}

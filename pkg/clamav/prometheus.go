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
	threadsIdle *prometheus.Desc
	threadsMax  *prometheus.Desc
}

func NewCollector(config cfg.Config) *Collector {
	return &Collector{
		config:      config,
		status:      prometheus.NewDesc("clamav_status", "Shows UP Status", nil, nil),
		threadsLive: prometheus.NewDesc("clamav_threads_live", "Shows live threads", nil, nil),
		threadsIdle: prometheus.NewDesc("clamav_threads_idle", "Shows idle threads", nil, nil),
		threadsMax:  prometheus.NewDesc("clamav_threads_max", "Shows max threads", nil, nil),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.status
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	address := fmt.Sprintf("%s:%d", collector.config.ClamAVAddress, collector.config.ClamAVPort)
	ch <- prometheus.MustNewConstMetric(collector.status, prometheus.CounterValue, ping(address))

	stats := stats(address)
	ch <- prometheus.MustNewConstMetric(collector.threadsLive, prometheus.CounterValue, stats.Threads.Live)
	ch <- prometheus.MustNewConstMetric(collector.threadsIdle, prometheus.CounterValue, stats.Threads.Idle)
	ch <- prometheus.MustNewConstMetric(collector.threadsMax, prometheus.CounterValue, stats.Threads.Max)
}

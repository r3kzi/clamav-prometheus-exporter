package clamav

import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/cfg"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	"regexp"
)

type Collector struct {
	config      cfg.Config
	status      *prometheus.Desc
	threadsLive *prometheus.Desc
	threadsIdle *prometheus.Desc
	threadsMax  *prometheus.Desc
	memHeap     *prometheus.Desc
	memMmap     *prometheus.Desc
	memUsed     *prometheus.Desc
}

var (
	regex = regexp.MustCompile("(live|idle|max|heap|mmap|\\bused)\\s([0-9.]+)[MG]*")
)

func NewCollector(config cfg.Config) *Collector {
	return &Collector{
		config:      config,
		status:      prometheus.NewDesc("clamav_status", "Shows UP Status", nil, nil),
		threadsLive: prometheus.NewDesc("clamav_threads_live", "Shows live threads", nil, nil),
		threadsIdle: prometheus.NewDesc("clamav_threads_idle", "Shows idle threads", nil, nil),
		threadsMax:  prometheus.NewDesc("clamav_threads_max", "Shows max threads", nil, nil),
		memHeap:     prometheus.NewDesc("clamav_mem_heap", "Shows heap memory usage", nil, nil),
		memMmap:     prometheus.NewDesc("clamav_mem_mmap", "Shows mmap memory usage", nil, nil),
		memUsed:     prometheus.NewDesc("clamav_mem_used", "Shows used memory usage", nil, nil),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.status
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	address := fmt.Sprintf("%s:%d", collector.config.ClamAVAddress, collector.config.ClamAVPort)

	pong := dial(address, commands.PING)
	if bytes.Compare(pong, []byte{'P', 'O', 'N', 'G', '\n'}) == 0 {
		ch <- prometheus.MustNewConstMetric(collector.status, prometheus.CounterValue, 1)
	}

	stats := dial(address, commands.STATS)
	matches := regex.FindAllStringSubmatch(string(stats), -1)
	if len(matches) > 0 {
		ch <- prometheus.MustNewConstMetric(collector.threadsLive, prometheus.CounterValue, toFloat(matches[0][2]))
		ch <- prometheus.MustNewConstMetric(collector.threadsIdle, prometheus.CounterValue, toFloat(matches[1][2]))
		ch <- prometheus.MustNewConstMetric(collector.threadsMax, prometheus.CounterValue, toFloat(matches[2][2]))
		ch <- prometheus.MustNewConstMetric(collector.memHeap, prometheus.GaugeValue, toFloat(matches[3][2]))
		ch <- prometheus.MustNewConstMetric(collector.memMmap, prometheus.GaugeValue, toFloat(matches[4][2]))
		ch <- prometheus.MustNewConstMetric(collector.memUsed, prometheus.GaugeValue, toFloat(matches[5][2]))
	}
}

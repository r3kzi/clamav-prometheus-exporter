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
	//THREADS: live 1  idle 0 max 12
	threadsRegex = regexp.MustCompile("(live)\\s*([0-9]+)\\s*(idle)\\s*([0-9]+)\\s*(max)\\s*([0-9]+)")
	//MEMSTATS: heap 3.656M mmap 0.129M used 3.236M free 0.420M releasable 0.127M pools 1 pools_used 1089.550M pools_total 1089.585M
	memRegex = regexp.MustCompile("(heap)\\s*([0-9.]+)([MG]+)\\s*(mmap)\\s*([0-9.]+)([MG]+)\\s*(used)\\s*([0-9.]+)([MG]+)\\s*(free)\\s*([0-9.]+)([MG]+)\\s*(releasable)\\s*([0-9.]+)([MG]+)\\s*(pools)\\s*([0-9]+)\\s*(pools_used)\\s*([0-9.]+)([MG]+)\\s*(pools_total)\\s*([0-9.]+)([MG]+)")
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

	threads := threadsRegex.FindAllStringSubmatch(string(stats), 1)
	if len(threads) > 0 {
		ch <- prometheus.MustNewConstMetric(collector.threadsLive, prometheus.CounterValue, toFloat(threads[0][2]))
		ch <- prometheus.MustNewConstMetric(collector.threadsIdle, prometheus.CounterValue, toFloat(threads[0][4]))
		ch <- prometheus.MustNewConstMetric(collector.threadsMax, prometheus.CounterValue, toFloat(threads[0][6]))
	}

	mem := memRegex.FindAllStringSubmatch(string(stats), 1)
	if len(mem) > 0 {
		ch <- prometheus.MustNewConstMetric(collector.memHeap, prometheus.GaugeValue, toFloat(mem[0][2]))
		ch <- prometheus.MustNewConstMetric(collector.memMmap, prometheus.GaugeValue, toFloat(mem[0][5]))
		ch <- prometheus.MustNewConstMetric(collector.memUsed, prometheus.GaugeValue, toFloat(mem[0][8]))
	}
}

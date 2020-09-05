package clamav

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	"io/ioutil"
	"strings"
)

type PrometheusCollector struct {
	status *prometheus.Desc
}

func NewPrometheusCollector() *PrometheusCollector {
	return &PrometheusCollector{
		status: prometheus.NewDesc("clamav_status",
			"Shows UP Status",
			nil, nil,
		),
	}
}

func (collector *PrometheusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.status
}

func (collector *PrometheusCollector) Collect(ch chan<- prometheus.Metric) {
	conn := NewTCPClient("localhost", 3310)
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("%s", commands.PING)))
	resp, _ := ioutil.ReadAll(conn)
	if strings.TrimSpace(string(resp)) == "PONG" {
		ch <- prometheus.MustNewConstMetric(collector.status, prometheus.CounterValue, 1)
	} else {
		ch <- prometheus.MustNewConstMetric(collector.status, prometheus.CounterValue, 0)
	}
}

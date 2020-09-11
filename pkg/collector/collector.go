/*
Copyright 2020 Christian Niehoff.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/clamav"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	"regexp"
	"strconv"
)

//Collector satisfies prometheus.Collector interface
type Collector struct {
	client      clamav.Client
	up          *prometheus.Desc
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

//New creates a Collector struct
func New(client clamav.Client) *Collector {
	return &Collector{
		client:      client,
		up:          prometheus.NewDesc("clamav_up", "Shows UP Status", nil, nil),
		threadsLive: prometheus.NewDesc("clamav_threads_live", "Shows live threads", nil, nil),
		threadsIdle: prometheus.NewDesc("clamav_threads_idle", "Shows idle threads", nil, nil),
		threadsMax:  prometheus.NewDesc("clamav_threads_max", "Shows max threads", nil, nil),
		memHeap:     prometheus.NewDesc("clamav_mem_heap", "Shows heap memory usage", nil, nil),
		memMmap:     prometheus.NewDesc("clamav_mem_mmap", "Shows mmap memory usage", nil, nil),
		memUsed:     prometheus.NewDesc("clamav_mem_used", "Shows used memory usage", nil, nil),
	}
}

//Describe satisfies prometheus.Collector.Describe
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.up
}

//Collect satisfies prometheus.Collector.Collect
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	pong := collector.client.Dial(commands.PING)
	if bytes.Compare(pong, []byte{'P', 'O', 'N', 'G', '\n'}) == 0 {
		ch <- prometheus.MustNewConstMetric(collector.up, prometheus.CounterValue, 1)
	}

	float := func(s string) float64 {
		float, err := strconv.ParseFloat(s, 64)
		if err != nil {
			_ = fmt.Errorf("couldn't parse string to float: %s", err)
		}
		return float
	}

	stats := collector.client.Dial(commands.STATS)
	matches := regex.FindAllStringSubmatch(string(stats), -1)
	if len(matches) > 0 {
		ch <- prometheus.MustNewConstMetric(collector.threadsLive, prometheus.CounterValue, float(matches[0][2]))
		ch <- prometheus.MustNewConstMetric(collector.threadsIdle, prometheus.CounterValue, float(matches[1][2]))
		ch <- prometheus.MustNewConstMetric(collector.threadsMax, prometheus.CounterValue, float(matches[2][2]))
		ch <- prometheus.MustNewConstMetric(collector.memHeap, prometheus.GaugeValue, float(matches[3][2]))
		ch <- prometheus.MustNewConstMetric(collector.memMmap, prometheus.GaugeValue, float(matches[4][2]))
		ch <- prometheus.MustNewConstMetric(collector.memUsed, prometheus.GaugeValue, float(matches[5][2]))
	}
}

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
	"math"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/clamav"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/commands"
	log "github.com/sirupsen/logrus"
)

//Collector satisfies prometheus.Collector interface
type Collector struct {
	client      clamav.Client
	up          *prometheus.Desc
	threadsLive *prometheus.Desc
	threadsIdle *prometheus.Desc
	threadsMax  *prometheus.Desc
	queue       *prometheus.Desc
	memHeap     *prometheus.Desc
	memMmap     *prometheus.Desc
	memUsed     *prometheus.Desc
	poolsUsed   *prometheus.Desc
	poolsTotal  *prometheus.Desc
	buildInfo   *prometheus.Desc
}

//New creates a Collector struct
func New(client clamav.Client) *Collector {
	return &Collector{
		client:      client,
		up:          prometheus.NewDesc("clamav_up", "Shows UP Status", nil, nil),
		threadsLive: prometheus.NewDesc("clamav_threads_live", "Shows live threads", nil, nil),
		threadsIdle: prometheus.NewDesc("clamav_threads_idle", "Shows idle threads", nil, nil),
		threadsMax:  prometheus.NewDesc("clamav_threads_max", "Shows max threads", nil, nil),
		queue:       prometheus.NewDesc("clamav_queue_length", "Shows queued items", nil, nil),
		memHeap:     prometheus.NewDesc("clamav_mem_heap_bytes", "Shows heap memory usage in bytes", nil, nil),
		memMmap:     prometheus.NewDesc("clamav_mem_mmap_bytes", "Shows mmap memory usage in bytes", nil, nil),
		memUsed:     prometheus.NewDesc("clamav_mem_used_bytes", "Shows used memory in bytes", nil, nil),
		poolsUsed:   prometheus.NewDesc("clamav_pools_used_bytes", "Shows memory used by memory pool allocator for the signature database in bytes", nil, nil),
		poolsTotal:  prometheus.NewDesc("clamav_pools_total_bytes", "Shows total memory allocated by memory pool allocator for the signature database in bytes", nil, nil),
		buildInfo:   prometheus.NewDesc("clamav_build_info", "Shows ClamAV Build Info", []string{"clamav_version", "database_version"}, nil),
	}
}

//Describe satisfies prometheus.Collector.Describe
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.up
	ch <- collector.threadsLive
	ch <- collector.threadsIdle
	ch <- collector.threadsMax
	ch <- collector.queue
	ch <- collector.memHeap
	ch <- collector.memMmap
	ch <- collector.memUsed
	ch <- collector.poolsUsed
	ch <- collector.poolsTotal
	ch <- collector.buildInfo
}

//Collect satisfies prometheus.Collector.Collect
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	pong := collector.client.Dial(commands.PING)
	if bytes.Equal(pong, []byte{'P', 'O', 'N', 'G', '\n'}) {
		ch <- prometheus.MustNewConstMetric(collector.up, prometheus.GaugeValue, 1)
	} else {
		ch <- prometheus.MustNewConstMetric(collector.up, prometheus.GaugeValue, 0)
	}

	float := func(s string) float64 {
		float, err := strconv.ParseFloat(s, 64)
		if err != nil {
			float = math.NaN()
		}
		return float
	}

	stats := collector.client.Dial(commands.STATS)
	idle, err := regexp.MatchString("IDLE", string(stats))
	if err != nil {
		log.Errorf("error searching IDLE field in stats %s: %s", idle, err)
		return
	}
	regex := regexp.MustCompile(`([0-9.]+|N/A)`)
	matches := regex.FindAllStringSubmatch(string(stats), -1)
	if len(matches) > 0 && idle == false {
		ch <- prometheus.MustNewConstMetric(collector.threadsLive, prometheus.GaugeValue, float(matches[1][1]))
		ch <- prometheus.MustNewConstMetric(collector.threadsIdle, prometheus.GaugeValue, float(matches[2][1]))
		ch <- prometheus.MustNewConstMetric(collector.threadsMax, prometheus.GaugeValue, float(matches[3][1]))
		ch <- prometheus.MustNewConstMetric(collector.queue, prometheus.GaugeValue, float(matches[5][1]))
		ch <- prometheus.MustNewConstMetric(collector.memHeap, prometheus.GaugeValue, float(matches[7][1])*1024)
		ch <- prometheus.MustNewConstMetric(collector.memMmap, prometheus.GaugeValue, float(matches[8][1])*1024)
		ch <- prometheus.MustNewConstMetric(collector.memUsed, prometheus.GaugeValue, float(matches[9][1])*1024)
		ch <- prometheus.MustNewConstMetric(collector.poolsUsed, prometheus.GaugeValue, float(matches[13][1])*1024)
		ch <- prometheus.MustNewConstMetric(collector.poolsTotal, prometheus.GaugeValue, float(matches[14][1])*1024)
	}

	if len(matches) > 0 && idle == true {
		ch <- prometheus.MustNewConstMetric(collector.threadsLive, prometheus.GaugeValue, float(matches[1][1]))
		ch <- prometheus.MustNewConstMetric(collector.threadsIdle, prometheus.GaugeValue, float(matches[2][1]))
		ch <- prometheus.MustNewConstMetric(collector.threadsMax, prometheus.GaugeValue, float(matches[3][1]))
		ch <- prometheus.MustNewConstMetric(collector.queue, prometheus.GaugeValue, float(matches[5][1]))
		ch <- prometheus.MustNewConstMetric(collector.memHeap, prometheus.GaugeValue, float(matches[8][1])*1024)
		ch <- prometheus.MustNewConstMetric(collector.memMmap, prometheus.GaugeValue, float(matches[9][1])*1024)
		ch <- prometheus.MustNewConstMetric(collector.memUsed, prometheus.GaugeValue, float(matches[10][1])*1024)
		ch <- prometheus.MustNewConstMetric(collector.poolsUsed, prometheus.GaugeValue, float(matches[14][1])*1024)
		ch <- prometheus.MustNewConstMetric(collector.poolsTotal, prometheus.GaugeValue, float(matches[15][1])*1024)
	}

	version := collector.client.Dial(commands.VERSION)
	regex = regexp.MustCompile(`((ClamAV)+\s([0-9.]*)/([0-9.]*))`)
	matches = regex.FindAllStringSubmatch(string(version), -1)
	if len(matches) > 0 {
		ch <- prometheus.MustNewConstMetric(collector.buildInfo, prometheus.GaugeValue, 1, matches[0][3], matches[0][4])
	}

}

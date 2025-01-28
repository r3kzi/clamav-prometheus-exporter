<p align="left"><img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/9e5f19f595edf1bb1a51cb49e4eac9f935c1ec18.png" alt="Logo" height="200"></p>

# ClamAV Prometheus Exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/r3kzi/clamav-prometheus-exporter)](https://goreportcard.com/report/github.com/r3kzi/clamav-prometheus-exporter)
[![Apache V2 License](https://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/r3kzi/clamav-prometheus-exporter/blob/master/LICENSE)

Exports metrics from [ClamAV](https://www.clamav.net/) as Prometheus metrics.

## Currently exposed metrics

- ClamAVUp
- ClamAVThreadsLive
- ClamAVThreadsIdle
- ClamAVThreadsMax
- ClamAVQueue
- ClamAVMemHeap
- ClamAVMemMMap
- ClamAVMemUsed
- ClamAVBuildInfo
- ClamAVDatabaseAge

```
# HELP clamav_build_info Shows ClamAV Build Info
# TYPE clamav_build_info gauge
clamav_build_info{clamav_version="0.102.4",database_version="26091"} 1
# HELP clamav_mem_heap_bytes Shows heap memory usage in bytes
# TYPE clamav_mem_heap_bytes gauge
clamav_mem_heap_bytes 1.090783104e+06
# HELP clamav_mem_mmap_bytes Shows mmap memory usage in bytes
# TYPE clamav_mem_mmap_bytes gauge
clamav_mem_mmap_bytes 1.076747264e+06
# HELP clamav_mem_used_bytes Shows used memory in bytes
# TYPE clamav_mem_used_bytes gauge
clamav_mem_used_bytes 1.076783104e+06
# HELP clamav_pools_total_bytes Shows total memory allocated by memory pool allocator for the signature database in bytes
# TYPE clamav_pools_total_bytes gauge
clamav_pools_total_bytes 1.076783104e+06
# HELP clamav_pools_used_bytes Shows memory used by memory pool allocator for the signature database in bytes
# TYPE clamav_pools_used_bytes gauge
clamav_pools_used_bytes 1.076747264e+06
# HELP clamav_queue_length Shows queued items
# TYPE clamav_queue_length gauge
clamav_queue_length 0
# HELP clamav_threads_idle Shows idle threads
# TYPE clamav_threads_idle gauge
clamav_threads_idle 0
# HELP clamav_threads_live Shows live threads
# TYPE clamav_threads_live gauge
clamav_threads_live 1
# HELP clamav_threads_max Shows max threads
# TYPE clamav_threads_max gauge
clamav_threads_max 10
# HELP clamav_up Shows UP Status
# TYPE clamav_up gauge
clamav_up 1
# HELP clamav_database_age Shows ClamAV signature database age in seconds
# TYPE clamav_database_age gauge
clamav_database_age 447408.4671055
```

## Installation

ClamAV Prometheus Exporter requires a
[supported release of Go](https://golang.org/doc/devel/release.html#policy).

```shell script
$ go get -u github.com/r3kzi/clamav-prometheus-exporter
```

To find out where `clamav-prometheus-exporter` was installed you can run `$ go list -f {{.Target}} github.com/r3kzi/clamav-prometheus-exporter`.

For `clamav-prometheus-exporter` to be used globally add that directory to the `$PATH` environment setting.

You could also build the binary yourself running:
```shell script
$ CGO_ENABLED=0 && go build -installsuffix 'static' -o clamav-prometheus-exporter .
```

## Flags

[ClamAV](https://www.clamav.net/) server to connect to:

```shell script
Usage of clamav-prometheus-exporter:
  -clamav-address string
    	ClamAV address to use (default "localhost")
  -clamav-port int
    	ClamAV port to use (default 3310)
  -log-level string
    	Log level (trace, debug, info, warn, error, fatal, panic) (default "info")
  -network string
    	Network mode to use, typically tcp or unix (socket) (default "tcp")
```

## Prometheus config

Just scrape this, e.g.:

```yaml
scrape_configs:
  - job_name: 'clamav-prometheus-exporter'
    static_configs:
      - targets: ['localhost:9810']
```

## Contributing

Pull requests are welcome.

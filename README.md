<p align="left"><img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/9e5f19f595edf1bb1a51cb49e4eac9f935c1ec18.png" alt="Logo" height="200"></p>

# ClamAV Prometheus Exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/r3kzi/clamav-prometheus-exporter)](https://goreportcard.com/report/github.com/r3kzi/clamav-prometheus-exporter)
[![Apache V2 License](https://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/r3kzi/clamav-prometheus-exporter/blob/master/LICENSE)

Exports metrics from [ClamAV](https://www.clamav.net/) as Prometheus metrics.

## Currently exposed metrics

- ClamAVBuildInfo
- ClamAVDatabaseAge
- ClamAVMemHeap
- ClamAVMemMmap
- ClamAVMemUsed
- ClamAVPoolsTotal
- ClamAVPoolsUsed
- ClamAVQueue
- ClamAVThreadsIdle
- ClamAVThreadsLive
- ClamAVThreadsMax
- ClamAVUp
- ClamAVVirusesDetected (requires `VirusEvent` in `clamd.conf`, see [Counting detected viruses](#counting-detected-viruses))

```
# HELP clamav_build_info Shows ClamAV Build Info
# TYPE clamav_build_info gauge
clamav_build_info{clamav_version="0.102.4",database_version="26091"} 1
# HELP clamav_database_age Shows ClamAV signature database age in seconds
# TYPE clamav_database_age gauge
clamav_database_age 447408.4671055
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
# HELP clamav_viruses_detected_total Shows the number of detected viruses reported by the VirusEvent directive of clamd
# TYPE clamav_viruses_detected_total counter
clamav_viruses_detected_total 0
```

## Installation

ClamAV Prometheus Exporter requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).

```shell
$ go get -u github.com/r3kzi/clamav-prometheus-exporter
```

To find out where `clamav-prometheus-exporter` was installed you can run `$ go list -f {{.Target}} github.com/r3kzi/clamav-prometheus-exporter`.

For `clamav-prometheus-exporter` to be used globally add that directory to the `$PATH` environment setting.

You could also build the binary yourself running:

```shell
$ make build
```

## Flags

[ClamAV](https://www.clamav.net/) server to connect to:

```shell
Usage of clamav-prometheus-exporter:
  -clamav-address string
      ClamAV address to use (default "localhost")
  -clamav-port int
      ClamAV port to use (default 3310)
  -log-level string
      Set the level of logging. (options: trace, debug, info, warn, error, fatal, panic) (default "info")
  -network string
      Network mode to use, typically tcp or unix (socket) (default "tcp")
  -virus-event-listen-address string
      Separate address to serve /virus-event on, e.g. 127.0.0.1:9811 (served on the metrics port if empty)
```

## Counting detected viruses

`clamd` does not expose the number of detected viruses through its socket commands. It can run a command
for every detection though (`VirusEvent`). The exporter counts every `POST` request to `/virus-event`, so
`clamd` has to send one per detection.

```
# HELP clamav_viruses_detected_total Shows the number of detected viruses reported by the VirusEvent directive of clamd
# TYPE clamav_viruses_detected_total counter
clamav_viruses_detected_total 2
```

### Required clamd configuration

Add one `VirusEvent` line to `clamd.conf` and restart `clamd`.

The command only has to send an empty `POST` request to the exporter. This line works with
`curl`, `wget` (GNU or BusyBox) and, if neither is installed, with `bash` or `nc`:

```
VirusEvent if command -v curl >/dev/null; then curl -s -o /dev/null -m 2 -X POST http://127.0.0.1:9810/virus-event; elif command -v wget >/dev/null; then wget -q -O /dev/null -T 2 --post-data= http://127.0.0.1:9810/virus-event; else timeout 2 bash -c 'printf "POST /virus-event HTTP/1.0\r\n\r\n" >/dev/tcp/127.0.0.1/9810'; fi
```

It must stay on a single line. If you know which tool is available, you can use the matching short form
instead (check with `command -v curl wget bash nc` in the environment of `clamd`):

```
# curl
VirusEvent curl -s -o /dev/null -m 2 -X POST http://127.0.0.1:9810/virus-event

# wget (GNU or BusyBox)
VirusEvent wget -q -O /dev/null -T 2 --post-data= http://127.0.0.1:9810/virus-event

# neither curl nor wget: bash writes the HTTP request to the TCP socket itself
VirusEvent timeout 2 bash -c 'printf "POST /virus-event HTTP/1.0\r\n\r\n" >/dev/tcp/127.0.0.1/9810'

# neither curl, wget nor bash: nc (netcat)
VirusEvent printf 'POST /virus-event HTTP/1.0\r\n\r\n' | nc -w 2 127.0.0.1 9810 >/dev/null
```

- Replace `127.0.0.1` and `9810` with the address and port of the exporter as seen from `clamd`, in every
  place of the line (the `bash` and `nc` forms take host and port separately). If the exporter runs with
  `-virus-event-listen-address`, use that address and port instead of the metrics port.
- `clamd` accepts only one `VirusEvent` directive. If you already use one, call both commands from a
  wrapper script.
- Keep the timeouts (`-m 2`, `-T 2`, `timeout 2`, `-w 2`), so that an unreachable exporter cannot hold up `clamd`.
- The request needs no body. `clamd` provides the signature name and the file name to the command
  (`%v`, `CLAM_VIRUSEVENT_VIRUSNAME`, `CLAM_VIRUSEVENT_FILENAME`), but the exporter does not use them.


### Things to know

- `VirusEvent` is a `clamd` feature. Scans with `clamscan` (without the daemon) are not counted.
- The endpoint has no authentication. Everyone who can reach it can increment the counter (nothing else).
  By default it is served on the metrics port. To restrict access, serve it on a separate address with
  `-virus-event-listen-address`: use `127.0.0.1:9811` if `clamd` runs next to the exporter (same host or pod),
  otherwise use e.g. `:9811` and only allow `clamd` to reach that port (firewall, Kubernetes NetworkPolicy).

## Prometheus config

Just scrape this, e.g.:

```yaml
scrape_configs:
  - job_name: 'clamav-prometheus-exporter'
    static_configs:
      - targets: ['localhost:9810']
```

## Release

For a new version of the application, bump version in the [VERSION](./VERSION) file.

>[!NOTE]
>The version in the file must follow the [Semantic Version](https://semver.org/) pattern: `X.Y.Z`.

This will trigger a  [GitHub Action](.github/workflows/docker-build.yml) that will be responsible
for creating a new image on [DockerHub](https://hub.docker.com/r/rekzi/clamav-prometheus-exporter)
in the pattern: `rekzi/clamav-prometheus-exporter:X.Y.Z`

In addition, a Release will be created in the repository in addition to the version patch increment,
following the pattern: `X.Y.Z-snapshot`

## Contributing

Pull requests are welcome.

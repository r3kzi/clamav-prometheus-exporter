# ClamAV Prometheus Exporter

## Currently exposed metrics

- ClamAVStatus
- ClamAVThreadsLive
- ClamAVThreadsIdle
- ClamAVThreadsMax

``` 
# HELP clamav_status Shows UP Status
# TYPE clamav_status counter
clamav_status 1
# HELP clamav_threads_idle Shows idle threads
# TYPE clamav_threads_idle counter
clamav_threads_idle 0
# HELP clamav_threads_live Shows live threads
# TYPE clamav_threads_live counter
clamav_threads_live 1
# HELP clamav_threads_max Shows max threads
# TYPE clamav_threads_max counter
clamav_threads_max 12
```
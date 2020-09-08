# ClamAV Prometheus Exporter

## Currently exposed metrics

- ClamAVStatus
- ClamAVThreadsLive
- ClamAVThreadsIdle
- ClamAVThreadsMax
- ClamAVMemHeap
- ClamAVMemMMap
- ClamAVMemUsed

``` 
# HELP clamav_mem_heap Shows heap memory usage
# TYPE clamav_mem_heap gauge
clamav_mem_heap 3.656
# HELP clamav_mem_mmap Shows mmap memory usage
# TYPE clamav_mem_mmap gauge
clamav_mem_mmap 0.129
# HELP clamav_mem_used Shows used memory usage
# TYPE clamav_mem_used gauge
clamav_mem_used 3.237
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
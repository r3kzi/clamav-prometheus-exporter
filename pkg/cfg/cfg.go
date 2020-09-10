package cfg

// Config is the root config for ClamAV Prometheus Exporter
type Config struct {
	ClamAVAddress string
	ClamAVPort    int
}

// Package virusevent counts the viruses that clamd reports through its
// VirusEvent directive.
//
// clamd does not expose the number of detected viruses through any of its
// socket commands. It can however run a command for every detection
// (VirusEvent in clamd.conf). That command sends a POST request to the
// Handler of this package, which increments a counter.
package virusevent

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Path is the default path the Handler is served on.
const Path = "/virus-event"

// Handler is a http.Handler that counts one detected virus per POST request.
// It also satisfies prometheus.Collector.
type Handler struct {
	detected prometheus.Counter
}

// New creates a Handler.
func New() *Handler {
	return &Handler{
		detected: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "clamav_viruses_detected_total",
			Help: "Shows the number of detected viruses reported by the VirusEvent directive of clamd",
		}),
	}
}

// Describe satisfies prometheus.Collector.Describe
func (h *Handler) Describe(ch chan<- *prometheus.Desc) {
	h.detected.Describe(ch)
}

// Collect satisfies prometheus.Collector.Collect
func (h *Handler) Collect(ch chan<- prometheus.Metric) {
	h.detected.Collect(ch)
}

// ServeHTTP satisfies http.Handler. Only POST requests are counted, so that
// crawlers and browsers do not increment the counter by accident.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.detected.Inc()
	log.Debug("virus event received from ", r.RemoteAddr)
	w.WriteHeader(http.StatusNoContent)
}

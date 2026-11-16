package virusevent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	h := New()
	assert.Equal(t, float64(0), testutil.ToFloat64(h))

	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodPut} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(method, Path, nil))
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		assert.Equal(t, http.MethodPost, rec.Header().Get("Allow"))
	}
	assert.Equal(t, float64(0), testutil.ToFloat64(h))

	for i := 1; i <= 3; i++ {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, Path, strings.NewReader("")))
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, float64(i), testutil.ToFloat64(h))
	}
}

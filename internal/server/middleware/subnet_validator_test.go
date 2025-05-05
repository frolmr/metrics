package middleware

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTrustedSubnet(t *testing.T) {
	tests := []struct {
		name           string
		trustedSubnet  *net.IPNet
		realIP         string
		expectedStatus int
	}{
		{
			name:           "no trusted subnet - allow all",
			trustedSubnet:  nil,
			realIP:         "192.168.1.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "IP in trusted subnet",
			trustedSubnet:  parseCIDR(t, "192.168.1.0/24"),
			realIP:         "192.168.1.10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "IP not in trusted subnet",
			trustedSubnet:  parseCIDR(t, "192.168.2.0/24"),
			realIP:         "10.0.0.1",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid IP format",
			trustedSubnet:  parseCIDR(t, "192.168.3.0/24"),
			realIP:         "not.an.ip",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "empty X-Real-IP header",
			trustedSubnet:  parseCIDR(t, "192.168.4.0/24"),
			realIP:         "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/", nil)
			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}

			rec := httptest.NewRecorder()

			middleware := WithTrustedSubnet(tt.trustedSubnet)
			middleware(handler).ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func parseCIDR(t *testing.T, cidr string) *net.IPNet {
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Fatalf("Failed to parse CIDR %s: %v", cidr, err)
	}
	return subnet
}

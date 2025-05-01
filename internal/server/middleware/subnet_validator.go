package middleware

import (
	"net"
	"net/http"
)

func WithTrustedSubnet(trustedSubnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if trustedSubnet == nil {
				next.ServeHTTP(res, req)
				return
			}

			reqIP := net.ParseIP(req.Header.Get("X-Real-IP"))
			if reqIP == nil {
				res.WriteHeader(http.StatusForbidden)
				return
			}

			if !trustedSubnet.Contains(reqIP) {
				res.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}

package auth

import (
	"log"
	"errors"
	"strings"
	"net/http"
)

func AuthClientMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("AuthClientMiddleware()")
		commonName, err := extractCNFromRequest(r)
		if err == nil {
			log.Printf("ClientCN: %s\n", commonName)
			next.ServeHTTP(w, r)
		}
	})
}

func AuthManagementMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("AuthManagementMiddleware()")
		commonName, err := extractCNFromRequest(r)
		if err == nil {
			log.Printf("ClientCN: %s\n", commonName)
			if strings.HasPrefix(commonName, "sysop.") {
				next.ServeHTTP(w, r)
			}
		}
	})
}

func AuthPrivateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("AuthPrivateMiddleware()")
		commonName, err := extractCNFromRequest(r)
		if err == nil {
			log.Printf("ClientCN: %s\n", commonName)
			if strings.HasPrefix(commonName, "sysop.") || strings.HasPrefix(commonName, "orchestrator.") || strings.HasPrefix(commonName, "authorization.") {
				next.ServeHTTP(w, r)
			}
		}
	})
}

// helpers

func extractCNFromRequest(r *http.Request) (string, error) {
	if r.TLS != nil && len(r.TLS.VerifiedChains) > 0 && len(r.TLS.VerifiedChains[0]) > 0 {
		var commonName = r.TLS.VerifiedChains[0][0].Subject.CommonName
		return commonName, nil
	}

	return "", errors.New("Could not get CN")
}

package auth

import (
	"log"
  "fmt"
	"errors"
	"strings"
	"net/http"
)

func AuthClientMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("AuthClientMiddleware()")
		commonName, err := extractCNFromRequest(r)
		if err == nil {
			log.Printf("ClientCN: %s\n", commonName)
			next.ServeHTTP(w, r)
		} else {
      log.Println(err)
    }
	})
}

func AuthManagementMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("AuthManagementMiddleware()")
		commonName, err := extractCNFromRequest(r)
		if err == nil {
			fmt.Printf("ClientCN: %s\n", commonName)
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
  log.Printf("extractCNFromRequest()")

  //fmt.Printf("%p\n", r.TLS)
  fmt.Printf("%+v\n", r.TLS)
  fmt.Printf("%d\n", len(r.TLS.VerifiedChains))

	if r.TLS != nil && len(r.TLS.VerifiedChains) > 0 && len(r.TLS.VerifiedChains[0]) > 0 {
		var commonName = r.TLS.VerifiedChains[0][0].Subject.CommonName
		return commonName, nil
	}

/*
  log.Print("Certificate chain:")
	for i, cert := range state.PeerCertificates {
		subject := cert.Subject
		issuer := cert.Issuer
		log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
	log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")
*/
	return "", errors.New("Could not get CN")
}


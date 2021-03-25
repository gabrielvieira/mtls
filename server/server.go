package main

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	// io.WriteString(w, "Hello, world!\n")
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		fmt.Printf("Unauthenticated")
	} else {
		name := r.TLS.PeerCertificates[0].Subject.CommonName
		algo := r.TLS.PeerCertificates[0].SignatureAlgorithm.String()
		fmt.Printf("Hello, %s!\n", name)
		fmt.Printf("algo, %s!\n", algo)

		cert := r.TLS.PeerCertificates[0]
		fingerprint := md5.Sum(cert.Raw)

		var buf bytes.Buffer
		for i, f := range fingerprint {
			if i > 0 {
				fmt.Fprintf(&buf, ":")
			}
			fmt.Fprintf(&buf, "%02X", f)
		}
		fmt.Printf("Fingerprint for %s", buf.String())
	}

	fmt.Fprint(w, "ola gente deu certo")
}

func main() {
	// Set up a /hello resource handler
	http.HandleFunc("/cert", helloHandler)

	// Create a CA certificate pool and add cert.pem to it
	// caCert, err := ioutil.ReadFile("cert/server_cert.pem")
	caCertClient, err := ioutil.ReadFile("cert/client_cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)
	caCertPool.AppendCertsFromPEM(caCertClient)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS("cert/server_cert.pem", "cert/server_key.pem"))
}

package server

import (
	"golang.org/x/crypto/acme/autocert"
)

// NewLocalAutoCertManager will new a local AutoCert Manager
// @see https://github.com/golang/crypto/blob/master/acme/autocert/autocert.go
// @see https://www.captaincodeman.com/2017/05/07/automatic-https-with-free-ssl-certificates-using-go-lets-encrypt
func NewLocalAutoCertManager(domains []string) *autocert.Manager {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...), //your domain here
		Cache:      autocert.DirCache("certs"),         // folder for storing certificates
	}
	return &certManager
}

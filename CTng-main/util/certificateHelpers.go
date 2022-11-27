package util

import (
	"github.com/google/certificate-transparency-go/asn1"
	"github.com/google/certificate-transparency-go/x509"
)

var REVOKE_EXTENSION_ID = asn1.ObjectIdentifier{2, 5, 29, 20}

type Place struct {
	Vector int
	Index  int
}

func FindRevokePlace(cert *x509.Certificate) *Place {
	checkID := REVOKE_EXTENSION_ID
	var p Place
	for _, ext := range cert.Extensions {
		if ext.Id.Equal(checkID) {
			asn1.Unmarshal(ext.Value, &p)
			return &p
		}
	}
	return nil
}

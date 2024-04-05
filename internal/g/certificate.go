// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package g

import (
	"crypto/rsa"
	"crypto/x509"
)

type Certificate struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
}

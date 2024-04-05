// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-05, by liasica

package g

import (
	"crypto/rsa"
	"crypto/x509"
)

type Config struct {
	root struct {
		rootCA         *x509.Certificate
		rootPrivateKey *rsa.PrivateKey
	}
}

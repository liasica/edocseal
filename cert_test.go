// Copyright (C) conseal. 2024-present.
//
// Created at 2024-04-03, by liasica

package edocseal

import (
	"crypto/x509/pkix"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateCertificateAuthority(t *testing.T) {
	data, err := CreateCertificateAuthority(pkix.Name{
		Country:            []string{"中国"},                 // 国家
		Province:           []string{"北京市"},                // 省份
		Locality:           []string{"东城区"},                // 城市
		Organization:       []string{"张三"},                 // 证书持有者组织名称
		OrganizationalUnit: []string{"110101199008087218"}, // 证书持有者组织唯一标识
		StreetAddress:      nil,
		PostalCode:         nil,
		SerialNumber:       "",
		CommonName:         "Joash", // 证书持有者通用名，需保持唯一，否则验证会失败
		Names:              nil,
	}, time.Hour*24*365, 4096)
	require.NoError(t, err)
	t.Log(data)
}

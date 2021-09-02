package testutils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/bhrg3se/flahmingo-homework/services/auth/store"
)

const privateKey1 = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCMID6PEVTL1DQQdU+lGvHH1wTuhwRJcdsvHAasb6LwoLohkUFw
slBNe/R9OdETzR++KIQSHygkkUCZVHarBpZVZgpuXVAFIZf7D3ggPSAgvfahnfAr
aD94rGH4dq1cdQNJ8/6tByh1iOcM+CtlCDpAado+CTA0C5n6WVYZJrIb5wIDAQAB
AoGANcWF3CbdcF5dIFe0GXqOf45ukQ30wi5T2u8ZTICGeWpkIs932kRC9ojzmD5g
kgWsAa/Qhpe4MtgefvhflV+wvj29Jt9IiXKtT2w1ckzPlGrAFM3O0YFK6uNlckUE
0OLKUIC7/XMUYQvNGlCLU+5eY+eIfZr2/Iq0FWUAAufszSkCQQDsRGL8VqOvLCAx
IXygUF8mjgrPX1LyB4TGvoAv4abIa0/8iHq6alPycsouhF1CCiKDTN9uLYgfw0cW
hGCeHntLAkEAl9RE2M0SoE/r/rr7VdzinrjwaBBUDESNimhxe3E/Ng2lghv+3Od1
P/LcLGfiFUxnAmOAJEsFZpW5fWrhSDIEVQJBAJkuNWcY3QLjfvObnGtr8GTUztlo
GiTlDwaz1/QzPqjOOoumCLv31/lmxwKlurjoTTwHiQyr5IeDWgGVTfb2GEUCQQCW
EPCbeNZYLaaMeDPDgjdbul+j+7+XAIsFqoXABQb9Xi+gkhsuLHfvZRPJsEP92S5X
ZhFZzezgoExeci2JNiahAkAh5MWjK3A8tvv0b18pmu/5bYMEZ0Ku8SvIlWvx5OIa
/VCjt61mvRKD7XcXMSS/kD33+hlfVb5RupeQyBWT/kws
-----END RSA PRIVATE KEY-----`

func GetMockPrivateKey1() *rsa.PrivateKey {
	pemBytes, _ := pem.Decode([]byte(privateKey1))
	key, _ := x509.ParsePKCS1PrivateKey(pemBytes.Bytes)

	return key
}

const MockToken1 = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjYwNDgwMDAwMDAwMDAwMCwiaWF0IjoxNjMwNDc0MzM4LCJwaG9uZU51bWJlciI6InNvbWVOdW1iZXIifQ.JB8LcAsYiHb6ofA0zt7hg_x2lGcCR0AZJvxCPTxHCxmK82gsFaNgHHY3fTe5KYRyzRsDTlZTLwTeo-VLwps2lBju2QV7HebZ-hPo1kMl9J74SSesYBXFODKBms_8ywO7wSgqYCeNJqnM06n-Zw2hSSPVNte0111RhyBRogen8mA"

var MockUser1 = store.User{
	ID:          "someID",
	Name:        "Some User",
	IsVerified:  true,
	PhoneNumber: "someNumber",
}

var MockUser2 = store.User{
	ID:          "someID2",
	Name:        "Some User 2",
	IsVerified:  true,
	PhoneNumber: "21234567890",
}

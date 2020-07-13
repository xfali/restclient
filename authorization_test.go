// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import "testing"

func TestParseWWWAuthenticate(t *testing.T) {
	testWWWAuth := `WWW-Authenticate: Digest realm="testrealm@host.com",
        algorithm="md5",
        qop="auth,auth-int",
        nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093",
        opaque="5ccc069c403ebaf9f0171e9517f40e41"`
	auth := ParseWWWAuthenticate(testWWWAuth)
	t.Log(auth)
}

func TestDigestAuth_Refresh(t *testing.T) {
	testWWWAuth := `WWW-Authenticate: Digest realm="testrealm@host.com",
        algorithm="md5",
        qop="auth,auth-int",
        nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093",
        opaque="5ccc069c403ebaf9f0171e9517f40e41"`
	auth := ParseWWWAuthenticate(testWWWAuth)
	digestAuth := NewDigestAuth("user", "pw")
	err := digestAuth.Refresh("GET", "test.com", nil, auth)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(digestAuth)
}

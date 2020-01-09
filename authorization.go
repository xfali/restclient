/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package restclient

import (
    "crypto/md5"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "fmt"
    "hash"
    "io"
    "regexp"
    "strings"
    "time"
)

type BasicAuth struct {
    Username string
    Password string
}

func NewBasicAuth(username, password string) *BasicAuth {
    return &BasicAuth{
        Username: username,
        Password: password,
    }
}

type DigestAuth struct {
    Username    string
    Password    string

    realm       string
    nonce       string
    algorithm   string
    qop         string
    nonceCount  int
    clientNonce string
    opaque      string

    method string
    uri    string
    body   []byte
}

type WWWAuthenticate struct {
    Algorithm string
    Realm     string
    Qop       []string
    Nonce     string
    Opaque    string
}

func NewDigestAuth(username, password string) *DigestAuth {
    return &DigestAuth{
        Username: username,
        Password: password,
    }
}

func (da *DigestAuth) Refresh(method, uri string, body []byte, wwwAuth *WWWAuthenticate) error {
    da.opaque = wwwAuth.Opaque
    da.selectQop(wwwAuth.Qop)
    da.nonce = wwwAuth.Nonce
    da.realm = wwwAuth.Realm
    da.algorithm = wwwAuth.Algorithm
    if da.algorithm == "" {
        da.algorithm = "MD5"
    }
    da.nonceCount++
    s, err := da.hash(fmt.Sprintf("%d:%s", time.Now().UnixNano(), RandomId(6)))
    if err != nil {
        return err
    }
    da.clientNonce = s

    da.method = method
    da.uri = uri
    da.body = body

    return nil
}

func (da *DigestAuth) selectQop(qops []string) {
    if len(qops) == 0 {
        da.qop = ""
    }
    for _, v := range qops {
        if v == "auth" || v == "auth-int" {
            da.qop = v
            return
        }
    }
}

func (da *DigestAuth) a1() (string, error) {
    return da.hash(fmt.Sprintf("%s:%s:%s", da.Username, da.realm, da.Password))
}

func (da *DigestAuth) a2() (string, error) {
    if da.qop == "" || da.qop == "auth" {
        return da.hash(fmt.Sprintf("%s:%s", da.method, da.uri))
    } else if da.qop == "auth-int" {
        body, err := da.hash(string(da.body))
        if err != nil {
            return "", err
        }
        return da.hash(fmt.Sprintf("%s:%s:%s", da.method, da.uri, body))
    }
    return "", errors.New("A2 qop not support: " + da.qop)
}

func (da *DigestAuth) response() (string, error) {
    a1, err := da.a1()
    if err != nil {
        return "", err
    }
    a2, err := da.a2()
    if err != nil {
        return "", err
    }

    if da.qop == "" {
        return da.hash(fmt.Sprintf("%s:%s:%s", a1, da.nonce, a2))
    } else if da.qop == "auth" || da.qop == "auth-int" {
        return da.hash(fmt.Sprintf("%s:%s:%08x:%s:%s:%s", a1, da.nonce, da.nonceCount, da.clientNonce, da.qop, a2))
    }
    return "", errors.New("Response qop not support: " + da.qop)
}

func (da *DigestAuth) hash(s string) (string, error) {
    var h hash.Hash
    algorithm := strings.ToUpper(strings.TrimSpace(da.algorithm))
    if algorithm == "" || algorithm == "MD5" || algorithm == "MD5-SESS" {
        h = md5.New()
    } else if algorithm == "SHA-256" || algorithm == "SHA-256-SESS" {
        h = sha256.New()
    }

    if h != nil {
        _, err := io.WriteString(h, s)
        if err != nil {
            return "", err
        }
        return hex.EncodeToString(h.Sum(nil)), nil
    }

    return "", errors.New("algorithm not support " + algorithm)
}

func (da *DigestAuth) String() string {
    ret, _ := da.ToString()
    return ret
}

func (da *DigestAuth) ToString() (string, error) {
    buf := strings.Builder{}

    resp, err := da.response()
    buf.WriteString("Digest ")
    buf.WriteString(fmt.Sprintf(`username="%s",`, da.Username))
    buf.WriteString(fmt.Sprintf(`realm="%s",`, da.realm))
    buf.WriteString(fmt.Sprintf(`nonce="%s",`, da.nonce))
    buf.WriteString(fmt.Sprintf(`uri="%s",`, da.uri))
    buf.WriteString(fmt.Sprintf(`qop="%s",`, da.qop))
    buf.WriteString(fmt.Sprintf(`nc=%08x,`, da.nonceCount))
    buf.WriteString(fmt.Sprintf(`cnonce="%s",`, da.clientNonce))
    buf.WriteString(fmt.Sprintf(`response="%s",`, resp))
    buf.WriteString(fmt.Sprintf(`algorithm="%s",`, da.algorithm))
    buf.WriteString(fmt.Sprintf(`opaque="%s"`, da.opaque))

    return buf.String(), err
}

func RandomId(length int) string {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

func ParseWWWAuthenticate(s string) *WWWAuthenticate {
    var wwwAuth = WWWAuthenticate{}

    algorithmRgx := regexp.MustCompile(`algorithm="([^ ,]+)"`)
    algorithmMatch := algorithmRgx.FindStringSubmatch(s)
    if algorithmMatch != nil {
        wwwAuth.Algorithm = algorithmMatch[1]
    }

    realmRgx := regexp.MustCompile(`realm="(.+?)"`)
    realmMatch := realmRgx.FindStringSubmatch(s)
    if realmMatch != nil {
        wwwAuth.Realm = realmMatch[1]
    }

    qopRgx := regexp.MustCompile(`qop="(.+?)"`)
    qopMatch := qopRgx.FindStringSubmatch(s)
    if qopMatch != nil {
        wwwAuth.Qop = strings.Split(qopMatch[1], ",")
    }

    nonceRgx := regexp.MustCompile(`nonce="(.+?)"`)
    nonceMatch := nonceRgx.FindStringSubmatch(s)
    if nonceMatch != nil {
        wwwAuth.Nonce = nonceMatch[1]
    }

    opaqueRgx := regexp.MustCompile(`opaque="(.+?)"`)
    opaqueMatch := opaqueRgx.FindStringSubmatch(s)
    if opaqueMatch != nil {
        wwwAuth.Opaque = opaqueMatch[1]
    }

    return &wwwAuth
}

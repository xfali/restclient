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

type DigestAuth struct {
    Username    string
    Password    string
    Realm       string
    Nonce       string
    Algorithm   string
    Qop         string
    NonceCount  int
    ClientNonce string
    Opaque      string

    Method string
    Uri    string
    Body   []byte
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
    da.Opaque = wwwAuth.Opaque
    da.selectQop(wwwAuth.Qop)
    da.Nonce = wwwAuth.Nonce
    da.Realm = wwwAuth.Realm
    da.Algorithm = wwwAuth.Algorithm
    if da.Algorithm == "" {
        da.Algorithm = "MD5"
    }
    da.NonceCount++
    s, err := da.hash(fmt.Sprintf("%d:%s", time.Now().UnixNano(), RandomId(6)))
    if err != nil {
        return err
    }
    da.ClientNonce = s

    da.Method = method
    da.Uri = uri
    da.Body = body

    return nil
}

func (da *DigestAuth) selectQop(qops []string) {
    if len(qops) == 0 {
        da.Qop = ""
    }
    for _, v := range qops {
        if v == "auth" || v == "auth-int" {
            da.Qop = v
            return
        }
    }
}

func (da *DigestAuth) a1() (string, error) {
    return da.hash(fmt.Sprintf("%s:%s:%s", da.Username, da.Realm, da.Password))
}

func (da *DigestAuth) a2() (string, error) {
    if da.Qop == "" || da.Qop == "auth" {
        return da.hash(fmt.Sprintf("%s:%s", da.Method, da.Uri))
    } else if da.Qop == "auth-int" {
        body, err := da.hash(string(da.Body))
        if err != nil {
            return "", err
        }
        return da.hash(fmt.Sprintf("%s:%s:%s", da.Method, da.Uri, body))
    }
    return "", errors.New("A2 Qop not support: " + da.Qop)
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

    if da.Qop == "" {
        return da.hash(fmt.Sprintf("%s:%s:%s", a1, da.Nonce, a2))
    } else if da.Qop == "auth" || da.Qop == "auth-int" {
        return da.hash(fmt.Sprintf("%s:%s:%08x:%s:%s:%s", a1, da.Nonce, da.NonceCount, da.ClientNonce, da.Qop, a2))
    }
    return "", errors.New("Response Qop not support: " + da.Qop)
}

func (da *DigestAuth) hash(s string) (string, error) {
    var h hash.Hash
    algorithm := strings.ToUpper(strings.TrimSpace(da.Algorithm))
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

    var err error
    s, err1 := da.hash(da.Username)
    if err1 != nil {
        err = err1
    }
    resp, err2 := da.response()
    if err2 != nil {
        err = err2
    }

    buf.WriteString("Digest ")
    buf.WriteString(fmt.Sprintf(`username="%s",`, s))
    buf.WriteString(fmt.Sprintf(`realm="%s",`, da.Realm))
    buf.WriteString(fmt.Sprintf(`nonce="%s",`, da.Nonce))
    buf.WriteString(fmt.Sprintf(`uri="%s",`, da.Uri))
    buf.WriteString(fmt.Sprintf(`qop="%s",`, da.Qop))
    buf.WriteString(fmt.Sprintf(`nc="%d",`, da.NonceCount))
    buf.WriteString(fmt.Sprintf(`cnonce="%s",`, da.ClientNonce))
    buf.WriteString(fmt.Sprintf(`response="%s",`, resp))
    buf.WriteString(fmt.Sprintf(`algorithm="%s",`, da.Algorithm))
    buf.WriteString(fmt.Sprintf(`opaque="%s"`, da.Opaque))

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

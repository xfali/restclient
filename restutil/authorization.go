/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package restutil

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
}


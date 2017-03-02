// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Id string

type IdGenerator interface {
	Gen(req *http.Request) Id
	IsValid(id Id) bool
}

type Sha1Generator struct {
	hashKey string
}

func NewSha1Generator(hashKey string) *Sha1Generator {
	return &Sha1Generator{hashKey}
}

var _ IdGenerator = NewSha1Generator("test")

func GenRandKey(strength int) []byte {
	k := make([]byte, strength)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}

func (gen *Sha1Generator) Gen(req *http.Request) Id {
	bs := GenRandKey(24)
	if len(bs) == 0 {
		return Id("")
	}

	sig := fmt.Sprintf("%s%d%s", req.RemoteAddr, time.Now().UnixNano(), string(bs))

	h := hmac.New(sha1.New, []byte(gen.hashKey))
	fmt.Fprintf(h, "%s", sig)
	return Id(hex.EncodeToString(h.Sum(nil)))
}

func (gen *Sha1Generator) IsValid(id Id) bool {
	return len(id) == 40
}

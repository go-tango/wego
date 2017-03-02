// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"time"
)

type Session struct {
	id      Id
	maxAge  time.Duration
	manager *Sessions
	rw      http.ResponseWriter
}

func (session *Session) Id() Id {
	return session.id
}

func (session *Session) SetId(id Id) {
	session.id = id
}

func (session *Session) Get(key string) interface{} {
	return session.manager.Store.Get(session.id, key)
}

func (session *Session) Set(key string, value interface{}) {
	session.manager.Store.Set(session.id, key, value)
}

func (session *Session) Del(key string) bool {
	return session.manager.Store.Del(session.id, key)
}

func (session *Session) Release() {
	session.manager.Invalidate(session.rw, session)
}

func (session *Session) IsValid() bool {
	return session.manager.Generator.IsValid(session.id)
}

func (session *Session) SetMaxAge(maxAge time.Duration) {
	session.maxAge = maxAge
}

func (session *Session) Sessions() *Sessions {
	return session.manager
}

func (session *Session) SetSession(s *Session) {
	session.id = s.id
	session.maxAge = s.maxAge
	session.manager = s.manager
	session.rw = s.rw
}

type Sessioner interface {
	SetSession(*Session)
}

var _ Sessioner = &Session{}

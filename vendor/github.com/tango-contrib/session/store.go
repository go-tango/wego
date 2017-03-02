// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import "time"

type Store interface {
	Add(id Id) bool
	Exist(id Id) bool
	Clear(id Id) bool

	Get(id Id, key string) interface{}
	Set(id Id, key string, value interface{}) error
	Del(id Id, key string) bool

	SetMaxAge(maxAge time.Duration)
	SetIdMaxAge(id Id, maxAge time.Duration)

	Run() error
}

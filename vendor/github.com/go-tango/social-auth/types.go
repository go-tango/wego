// Copyright 2014 beego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
//
// Maintain by https://github.com/slene

package social

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/session"
)

// Interface of social Privider
type Provider interface {
	GetConfig() *Config
	GetType() SocialType
	GetName() string
	GetPath() string
	GetIndentify(*Token) (string, error)
	CanConnect(*Token, *UserSocial) (bool, error)
}

// Interface of social utils
type SocialAuther interface {
	IsUserLogin(*tango.Context, *session.Session) (int, bool)
	LoginUser(*tango.Context, *session.Session, int) (string, error)
}

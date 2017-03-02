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

package apps

import (
	"github.com/go-tango/social-auth"
)

type Weibo struct {
	BaseProvider
}

func (p *Weibo) GetType() social.SocialType {
	return social.SocialWeibo
}

func (p *Weibo) GetName() string {
	return "Weibo"
}

func (p *Weibo) GetPath() string {
	return "weibo"
}

func (p *Weibo) GetIndentify(tok *social.Token) (string, error) {
	return tok.GetExtra("uid"), nil
}

var _ social.Provider = new(Weibo)

func NewWeibo(clientId, secret string) *Weibo {
	p := new(Weibo)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "email"
	p.AuthURL = "https://api.weibo.com/oauth2/authorize"
	p.TokenURL = "https://api.weibo.com/oauth2/access_token"
	p.RedirectURL = social.DefaultAppUrl + "login/weibo/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}

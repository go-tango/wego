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

type Dropbox struct {
	BaseProvider
}

func (p *Dropbox) GetType() social.SocialType {
	return social.SocialDropbox
}

func (p *Dropbox) GetName() string {
	return "Dropbox"
}

func (p *Dropbox) GetPath() string {
	return "dropbox"
}

func (p *Dropbox) GetIndentify(tok *social.Token) (string, error) {
	return tok.GetExtra("uid"), nil
}

var _ social.Provider = new(Dropbox)

func NewDropbox(clientId, secret string) *Dropbox {
	p := new(Dropbox)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = ""
	p.AuthURL = "https://www.dropbox.com/1/oauth2/authorize"
	p.TokenURL = "https://api.dropbox.com/1/oauth2/token"
	p.RedirectURL = social.DefaultAppUrl + "login/dropbox/access"
	p.AccessType = ""
	p.ApprovalPrompt = ""
	return p
}

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
	"encoding/base64"

	"github.com/go-tango/social-auth"
)

type BaseProvider struct {
	App            social.Provider
	ClientId       string
	ClientSecret   string
	Scope          string
	AuthURL        string
	TokenURL       string
	RedirectURL    string
	AccessType     string
	ApprovalPrompt string
}

func (p *BaseProvider) getBasicAuth() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(p.ClientId+":"+p.ClientSecret))
}

func (p *BaseProvider) GetConfig() *social.Config {
	return &social.Config{
		ClientId:       p.ClientId,
		ClientSecret:   p.ClientSecret,
		Scope:          p.Scope,
		AuthURL:        p.AuthURL,
		TokenURL:       p.TokenURL,
		RedirectURL:    p.RedirectURL,
		AccessType:     p.AccessType,
		ApprovalPrompt: p.ApprovalPrompt,
	}
}

func (p *BaseProvider) CanConnect(tok *social.Token, userSocial *social.UserSocial) (bool, error) {
	identify, err := p.App.GetIndentify(tok)
	if err != nil {
		return false, err
	}

	has, err := social.ORM().Where("identify = ?", identify).And("type = ?", p.App.GetType()).Get(userSocial)
	if err != nil {
		return false, err
	}
	if !has {
		return true, nil
	}

	return false, nil
}

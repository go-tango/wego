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
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/httplib"

	"github.com/go-tango/social-auth"
)

type Github struct {
	BaseProvider
}

func (p *Github) GetType() social.SocialType {
	return social.SocialGithub
}

func (p *Github) GetName() string {
	return "Github"
}

func (p *Github) GetPath() string {
	return "github"
}

func (p *Github) GetIndentify(tok *social.Token) (string, error) {
	vals := make(map[string]interface{})

	uri := "https://api.github.com/user"
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)
	req.Header("Authorization", "Bearer "+tok.AccessToken)

	resp, err := req.Response()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()

	if err := decoder.Decode(&vals); err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%v", vals["message"])
	}

	if vals["id"] == nil {
		return "", nil
	}

	return fmt.Sprint(vals["id"]), nil
}

var _ social.Provider = new(Github)

func NewGithub(clientId, secret string) *Github {
	p := new(Github)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "user,public_repo"
	p.AuthURL = "https://github.com/login/oauth/authorize"
	p.TokenURL = "https://github.com/login/oauth/access_token"
	p.RedirectURL = social.DefaultAppUrl + "login/github/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}

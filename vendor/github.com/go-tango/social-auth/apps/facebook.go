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
	"net/url"

	"github.com/astaxie/beego/httplib"

	"github.com/go-tango/social-auth"
)

type Facebook struct {
	BaseProvider
}

func (p *Facebook) GetType() social.SocialType {
	return social.SocialFacebook
}

func (p *Facebook) GetName() string {
	return "Facebook"
}

func (p *Facebook) GetPath() string {
	return "facebook"
}

func (p *Facebook) GetIndentify(tok *social.Token) (string, error) {
	vals := make(map[string]interface{})

	uri := "https://graph.facebook.com/me?fields=id&access_token=" + url.QueryEscape(tok.AccessToken)
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

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

	if vals["error"] != nil {
		return "", fmt.Errorf("%v", vals["error"])
	}

	if vals["id"] == nil {
		return "", nil
	}

	return fmt.Sprint(vals["id"]), nil
}

var _ social.Provider = new(Facebook)

func NewFacebook(clientId, secret string) *Facebook {
	p := new(Facebook)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "email"
	p.AuthURL = "https://www.facebook.com/dialog/oauth"
	p.TokenURL = "https://graph.facebook.com/oauth/access_token"
	p.RedirectURL = social.DefaultAppUrl + "login/facebook/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}

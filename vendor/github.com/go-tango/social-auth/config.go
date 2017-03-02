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
// modified from https://code.google.com/p/goauth2/source/browse/oauth/oauth.go
//
// Maintain by https://github.com/slene

package social

import (
	"net/http"
	"net/url"
)

// default redirect url
var DefaultAppUrl = "http://127.0.0.1:8092/"

// default reansport
var DefaultTransport = http.DefaultTransport

// Config is the configuration of an OAuth consumer.
type Config struct {
	// ClientId is the OAuth client identifier used when communicating with
	// the configured OAuth provider.
	ClientId string

	// ClientSecret is the OAuth client secret used when communicating with
	// the configured OAuth provider.
	ClientSecret string

	// Scope identifies the level of access being requested. Multiple scope
	// values should be provided as a space-delimited string.
	Scope string

	// AuthURL is the URL the user will be directed to in order to grant
	// access.
	AuthURL string

	// TokenURL is the URL used to retrieve OAuth tokens.
	TokenURL string

	// RedirectURL is the URL to which the user will be returned after
	// granting (or denying) access.
	RedirectURL string

	// TokenCache allows tokens to be cached for subsequent requests.
	TokenCache Cache

	AccessType string // Optional, "online" (default) or "offline", no refresh token if "online"

	// ApprovalPrompt indicates whether the user should be
	// re-prompted for consent. If set to "auto" (default) the
	// user will be prompted only if they haven't previously
	// granted consent and the code can only be exchanged for an
	// access token.
	// If set to "force" the user will always be prompted, and the
	// code can be exchanged for a refresh token.
	ApprovalPrompt string
}

// AuthCodeURL returns a URL that the end-user should be redirected to,
// so that they may obtain an authorization code.
func (c *Config) AuthCodeURL(state string) string {
	url_, err := url.Parse(c.AuthURL)
	if err != nil {
		panic("AuthURL malformed: " + err.Error())
	}
	values := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientId},
		"redirect_uri":  {c.RedirectURL},
		"state":         {state},
	}

	if len(c.Scope) > 0 {
		values.Set("scope", c.Scope)
	}

	if len(c.AccessType) > 0 {
		values.Set("access_type", c.AccessType)
	}

	if len(c.ApprovalPrompt) > 0 {
		values.Set("approval_prompt", c.ApprovalPrompt)
	}

	q := values.Encode()
	if url_.RawQuery == "" {
		url_.RawQuery = q
	} else {
		url_.RawQuery += "&" + q
	}
	return url_.String()
}

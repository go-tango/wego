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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type OAuthError struct {
	prefix string
	msg    string
}

func (oe OAuthError) Error() string {
	return "OAuthError: " + oe.prefix + ": " + oe.msg
}

// Cache specifies the methods that implement a Token cache.
type Cache interface {
	Token() (*Token, error)
	PutToken(*Token) error
}

// Token contains an end-user's tokens.
// This is the data you must store to persist authentication.
type Token struct {
	AccessToken  string
	TokenType    string
	Expiry       time.Time // If zero the token has no (known) expiry time.
	RefreshToken string
	Extra        map[string]string // May be nil.
}

func (t *Token) Expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Before(time.Now())
}

func (t *Token) IsEmpty() bool {
	return len(t.AccessToken) == 0
}

func (t *Token) GetExtra(key string) string {
	if t.Extra == nil {
		return ""
	}
	return t.Extra[key]
}

// Transport implements http.RoundTripper. When configured with a valid
// Config and Token it can be used to make authenticated HTTP requests.
//
//	t := &oauth.Transport{config}
//      t.Exchange(code)
//      // t now contains a valid Token
//	r, _, err := t.Client().Get("http://example.org/url/requiring/auth")
//
// It will automatically refresh the Token if it can,
// updating the supplied Token in place.
type Transport struct {
	*Config
	*Token

	// Transport is the HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	// (It should never be an oauth.Transport.)
	Transport http.RoundTripper
}

// Client returns an *http.Client that makes OAuth-authenticated requests.
func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// Exchange takes a code and gets access Token from the remote server.
func (t *Transport) Exchange(code string) (*Token, error) {
	if t.Config == nil {
		return nil, OAuthError{"Exchange", "no Config supplied"}
	}

	// If the transport or the cache already has a token, it is
	// passed to `updateToken` to preserve existing refresh token.
	tok := t.Token
	if tok == nil && t.TokenCache != nil {
		tok, _ = t.TokenCache.Token()
	}
	if tok == nil {
		tok = new(Token)
	}

	values := url.Values{
		"grant_type":   {"authorization_code"},
		"redirect_uri": {t.RedirectURL},
		"code":         {code},
	}

	if len(t.Scope) > 0 {
		values.Set("scope", t.Scope)
	}

	err := t.updateToken(tok, values)
	if err != nil {
		return nil, err
	}
	t.Token = tok
	if t.TokenCache != nil {
		return tok, t.TokenCache.PutToken(tok)
	}
	return tok, nil
}

// RoundTrip executes a single HTTP transaction using the Transport's
// Token as authorization headers.
//
// This method will attempt to renew the Token if it has expired and may return
// an error related to that Token renewal before attempting the client request.
// If the Token cannot be renewed a non-nil os.Error value will be returned.
// If the Token is invalid callers should expect HTTP-level errors,
// as indicated by the Response's StatusCode.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Token == nil {
		if t.Config == nil {
			return nil, OAuthError{"RoundTrip", "no Config supplied"}
		}
		if t.TokenCache == nil {
			return nil, OAuthError{"RoundTrip", "no Token supplied"}
		}
		var err error
		t.Token, err = t.TokenCache.Token()
		if err != nil {
			return nil, err
		}
	}

	// Refresh the Token if it has expired.
	if t.Expired() {
		if err := t.Refresh(); err != nil {
			return nil, err
		}
	}

	// To set the Authorization header, we must make a copy of the Request
	// so that we don't modify the Request we were given.
	// This is required by the specification of http.RoundTripper.
	req = cloneRequest(req)
	req.Header.Set("Authorization", "Bearer "+t.AccessToken)
	req.Header.Set("Accept", "application/json")

	// Make the HTTP request.
	return t.transport().RoundTrip(req)
}

// Refresh renews the Transport's AccessToken using its RefreshToken.
func (t *Transport) Refresh() error {
	if t.Token == nil {
		return OAuthError{"Refresh", "no existing Token"}
	}
	if t.RefreshToken == "" {
		return OAuthError{"Refresh", "Token expired; no Refresh Token"}
	}
	if t.Config == nil {
		return OAuthError{"Refresh", "no Config supplied"}
	}

	err := t.updateToken(t.Token, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {t.RefreshToken},
	})
	if err != nil {
		return err
	}
	if t.TokenCache != nil {
		return t.TokenCache.PutToken(t.Token)
	}
	return nil
}

func (t *Transport) updateToken(tok *Token, v url.Values) error {
	v.Set("client_id", t.ClientId)
	v.Set("client_secret", t.ClientSecret)
	r, err := (&http.Client{Transport: t.transport()}).PostForm(t.TokenURL, v)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return OAuthError{"updateToken", r.Status}
	}

	parseValues := func(k string, v interface{}) {
		value := fmt.Sprint(v)
		switch k {
		case "access_token":
			tok.AccessToken = value
		case "token_type":
			tok.TokenType = value
		case "expires_in", "expires":
			d, _ := time.ParseDuration(value + "s")
			if d == 0 {
				tok.Expiry = time.Time{}
			} else {
				tok.Expiry = time.Now().Add(d)
			}
		case "refresh_token":
			// Don't overwrite `RefreshToken` with an empty value
			if len(value) != 0 {
				tok.RefreshToken = value
			}
		default:
			if tok.Extra == nil {
				tok.Extra = make(map[string]string)
			}
			tok.Extra[k] = value
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	body = bytes.TrimSpace(body)

	var content string

	if body[0] == '{' {
		content = "json"
	}

	switch content {
	case "json":
		vals := make(map[string]interface{})
		if err = json.Unmarshal(body, &vals); err != nil {
			return err
		}

		for key, value := range vals {
			parseValues(key, value)
		}
	default:
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return err
		}

		for key, _ := range vals {
			parseValues(key, vals.Get(key))
		}
	}

	return nil
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

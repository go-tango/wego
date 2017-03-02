// Copyright 2013 Beego Authors
// Copyright 2014 Unknwon
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

// Package captcha a middleware that provides captcha service for Macaron.
package captcha

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/Unknwon/com"
	"github.com/lunny/tango"
	"github.com/tango-contrib/cache"
)

var (
	defaultChars = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
)

// Captcha represents a captcha service.
type Captchas struct {
	Options
}

// NewCaptcha initializes and returns a captcha with given options.
func New(opts ...Options) *Captchas {
	opt := prepareOptions(opts)
	return &Captchas{
		Options: opt,
	}
}

// generate key string
func (c *Captchas) key(id string) string {
	return c.CachePrefix + id
}

// generate rand chars with default chars
func (c *Captchas) genRandChars() string {
	return string(com.RandomCreateBytes(c.ChallengeNums, defaultChars...))
}

func (c *Captchas) GenRandChars() string {
	return c.genRandChars()
}

// create a new captcha id
func (c *Captchas) CreateCaptcha() (string, error) {
	id := string(com.RandomCreateBytes(15))
	if err := c.Caches.Put(c.key(id), c.genRandChars(), c.Expiration); err != nil {
		return "", err
	}
	return id, nil
}

// verify from a request
func (c *Captchas) VerifyReq(req *http.Request) bool {
	req.ParseForm()
	return c.Verify(req.FormValue(c.FieldIdName), req.FormValue(c.FieldCaptchaName))
}

// direct verify id and challenge string
func (c *Captchas) Verify(id string, challenge string) bool {
	if len(challenge) == 0 || len(id) == 0 {
		return false
	}

	var chars string

	key := c.key(id)
	if v, ok := c.Caches.Get(key).(string); ok {
		chars = v
	} else {
		return false
	}

	defer c.Caches.Delete(key)

	if len(chars) != len(challenge) {
		return false
	}

	// verify challenge
	for i, c := range []byte(chars) {
		if c != challenge[i]-48 {
			return false
		}
	}

	return true
}

// tempalte func for output html
func (c *Captchas) CreateHtml() template.HTML {
	value, err := c.CreateCaptcha()
	if err != nil {
		panic(fmt.Errorf("fail to create captcha: %v", err))
	}
	return template.HTML(fmt.Sprintf(`<input type="hidden" name="%s" value="%s">
	<a class="captcha" href="javascript:">
		<img onclick="this.src=('%s%s%s.png?reload='+(new Date()).getTime())" class="captcha-img" src="%s%s%s.png">
	</a>`, c.FieldIdName, value, c.SubURL, c.URLPrefix, value, c.SubURL, c.URLPrefix, value))
}

type Options struct {
	Caches *cache.Caches
	// Suburl path. Default is empty.
	SubURL string
	// URL prefix of getting captcha pictures. Default is "/captcha/".
	URLPrefix string
	// Hidden input element ID. Default is "captcha_id".
	FieldIdName string
	// User input value element name in request form. Default is "captcha".
	FieldCaptchaName string
	// Challenge number. Default is 6.
	ChallengeNums int
	// Captcha image width. Default is 240.
	Width int
	// Captcha image height. Default is 80.
	Height int
	// Captcha expiration time in seconds. Default is 600.
	Expiration int64
	// Cache key prefix captcha characters. Default is "captcha_".
	CachePrefix string
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	if opt.Caches == nil {
		opt.Caches = cache.New(cache.Options{Adapter: "memory", Interval: 120})
	}

	opt.SubURL = strings.TrimSuffix(opt.SubURL, "/")

	// Defaults.
	if len(opt.URLPrefix) == 0 {
		opt.URLPrefix = "/captcha/"
	} else if opt.URLPrefix[len(opt.URLPrefix)-1] != '/' {
		opt.URLPrefix += "/"
	}
	if len(opt.FieldIdName) == 0 {
		opt.FieldIdName = "captcha_id"
	}
	if len(opt.FieldCaptchaName) == 0 {
		opt.FieldCaptchaName = "captcha"
	}
	if opt.ChallengeNums == 0 {
		opt.ChallengeNums = 6
	}
	if opt.Width == 0 {
		opt.Width = stdWidth
	}
	if opt.Height == 0 {
		opt.Height = stdHeight
	}
	if opt.Expiration == 0 {
		opt.Expiration = 600
	}
	if len(opt.CachePrefix) == 0 {
		opt.CachePrefix = "captcha_"
	}

	return opt
}

type Captcha struct {
	c   *Captchas
	req *http.Request
}

func (c *Captcha) SetCaptcha(cpt *Captchas, req *http.Request) {
	c.c = cpt
	c.req = req
}

func (c *Captcha) CreateHtml() template.HTML {
	return c.c.CreateHtml()
}

// verify from a request
func (c *Captcha) Verify() bool {
	return c.c.VerifyReq(c.req)
}

func (c *Captcha) CreateCaptcha() (string, error) {
	return c.c.CreateCaptcha()
}

func (c *Captcha) VerifyCaptcha(id string, challenge string) bool {
	return c.c.Verify(id, challenge)
}

type Captchaer interface {
	SetCaptcha(*Captchas, *http.Request)
}

// Captchaer is a middleware that maps a captcha.Captcha service into the Macaron handler chain.
// An single variadic captcha.Options struct can be optionally provided to configure.
// This should be register after cache.Cacher.
func (c *Captchas) Handle(ctx *tango.Context) {
	if !strings.HasPrefix(ctx.Req().RequestURI, c.URLPrefix) {
		if action := ctx.Action(); action != nil {
			if p, ok := action.(Captchaer); ok {
				p.SetCaptcha(c, ctx.Req())
			}
		}

		ctx.Next()
		return
	}

	var chars string
	id := path.Base(ctx.Req().RequestURI)
	if i := strings.Index(id, "."); i > -1 {
		id = id[:i]
	}
	key := c.key(id)

	// Reload captcha.
	if len(ctx.Req().FormValue("reload")) > 0 {
		chars = c.genRandChars()
		if err := c.Caches.Put(key, chars, c.Expiration); err != nil {
			ctx.WriteHeader(http.StatusInternalServerError)
			ctx.Write([]byte("captcha reload error"))
			panic(fmt.Errorf("fail to reload captcha: %v", err))
		}
	} else {
		if v, ok := c.Caches.Get(key).(string); ok {
			chars = v
		} else {
			ctx.WriteHeader(http.StatusNotFound)
			ctx.Write([]byte("captcha not found"))
			return
		}
	}

	if _, err := NewImage([]byte(chars), c.Width, c.Height).WriteTo(ctx.ResponseWriter); err != nil {
		panic(fmt.Errorf("fail to write captcha: %v", err))
	}
}

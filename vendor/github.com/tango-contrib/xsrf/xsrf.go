// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsrf

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-xweb/uuid"
	"github.com/lunny/tango"
)

const (
	XSRF_TAG string = "_xsrf"
)

type Xsrfer interface {
	CheckXsrf() bool
}

type NoCheck struct {
}

func (NoCheck) InitXsrfer(*tango.Context, time.Duration) {}

func (NoCheck) CheckXsrf() bool {
	return false
}

var _ Xsrfer = NoCheck{}

type XsrfChecker interface {
	SetXsrf(string, *tango.Context, time.Duration)
	AutoCheck() bool
}

type Checker struct {
	XsrfValue string
	ctx       *tango.Context
	timeout   time.Duration
}

func (Checker) CheckXsrf() bool {
	return true
}

func (c *Checker) SetXsrf(v string, ctx *tango.Context, timeout time.Duration) {
	c.XsrfValue = v
	c.ctx = ctx
	c.timeout = timeout
}

func (c *Checker) AutoCheck() bool {
	return true
}

func (c *Checker) XsrfFormHtml() template.HTML {
	return template.HTML(fmt.Sprintf(`<input type="hidden" name="%v" value="%v" />`,
		XSRF_TAG, c.XsrfValue))
}

func (c *Checker) Renew() {
	var val = uuid.NewRandom().String()
	var cookie = newCookie(XSRF_TAG, val, int64(c.timeout.Seconds()))
	c.ctx.Cookies().Del(XSRF_TAG)
	c.ctx.Cookies().Set(cookie)
	//c.ctx.Header().Set("Set-Cookie", cookie.String())
}

func (c *Checker) IsValid() bool {
	if c.ctx.Req().Method == "POST" {
		res, err := c.ctx.Req().Cookie(XSRF_TAG)
		formVal := c.ctx.Req().FormValue(XSRF_TAG)

		if err != nil || res.Value == "" || res.Value != formVal {
			return false
		}
	}

	return true
}

var _ XsrfChecker = &Checker{}

func New(timeout time.Duration) tango.HandlerFunc {
	return func(ctx *tango.Context) {
		var action interface{}
		if action = ctx.Action(); action == nil {
			ctx.Next()
			return
		}

		// if action implements check xsrf option and ask not check then return
		if checker, ok := action.(Xsrfer); ok && !checker.CheckXsrf() {
			ctx.Next()
			return
		}

		var val string = ""
		cookie, err := ctx.Req().Cookie(XSRF_TAG)
		if err != nil {
			val = uuid.NewRandom().String()
			cookie = newCookie(XSRF_TAG, val, int64(timeout.Seconds()))
			ctx.Cookies().Del(XSRF_TAG)
			ctx.Cookies().Set(cookie)
			//ctx.Header().Set("Set-Cookie", cookie.String())
		} else {
			val = cookie.Value
		}

		if c, ok := action.(XsrfChecker); ok {
			c.SetXsrf(val, ctx, timeout)

			if c.AutoCheck() {
				if ctx.Req().Method == "POST" {
					res, err := ctx.Req().Cookie(XSRF_TAG)
					formVal := ctx.Req().FormValue(XSRF_TAG)

					if err != nil || res.Value == "" || res.Value != formVal {
						ctx.Abort(http.StatusInternalServerError, "xsrf token error.")
						ctx.Error("xsrf token error.")
						return
					}
				}
			}
		}

		ctx.Next()
	}
}

// NewCookie is a helper method that returns a new http.Cookie object.
// Duration is specified in seconds. If the duration is zero, the cookie is permanent.
// This can be used in conjunction with ctx.SetCookie.
func newCookie(name string, value string, age int64) *http.Cookie {
	var utctime time.Time
	if age == 0 {
		// 2^31 - 1 seconds (roughly 2038)
		utctime = time.Unix(2147483647, 0)
	} else {
		utctime = time.Unix(time.Now().Unix()+age, 0)
	}
	return &http.Cookie{Name: name, Path: "/", Value: value, Expires: utctime}
}

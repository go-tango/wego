// Copyright 2015 wego authors
// Copyright 2013 wetalk authors
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

// An open source project for Gopher community.
package main

import (
	"time"

	"github.com/lunny/tango"
	"github.com/tango-contrib/debug"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"

	"github.com/go-tango/social-auth"
	"github.com/go-tango/wego/middlewares"
	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/routers"
	"github.com/go-tango/wego/routers/auth"
	"github.com/go-tango/wego/setting"

	_ "github.com/go-sql-driver/mysql"
)

func initTango(isprod bool) *tango.Tango {
	middlewares.Init()

	tg := tango.NewWithLog(setting.Log)

	if false {
		//if !isprod {
		tg.Use(debug.Debug(debug.Options{
			IgnorePrefix:     "/static",
			HideResponseBody: true,
			HideRequestBody:  true,
		}))
	}

	tg.Use(tango.ClassicHandlers...)

	sess := session.New(session.Options{
		MaxAge: time.Duration(setting.SessionCookieLifeTime),
	})
	tg.Use(
		tango.Static(tango.StaticOptions{
			RootPath: "./static",
			Prefix:   "static",
		}),
		tango.Static(tango.StaticOptions{
			RootPath: "./static_source",
			Prefix:   "static_source",
		}),
		sess,
		middlewares.Renders,
		setting.Captcha,
	)
	tg.Get("/favicon.ico", func(ctx *tango.Context) {
		ctx.ServeFile("./static/favicon.ico")
	})
	if setting.EnableXSRF {
		tg.Use(xsrf.New(time.Duration(setting.SessionCookieLifeTime)))
	}
	tg.Use(flash.Flashes(sess), events.Events())
	return tg
}

func main() {
	// init configs
	setting.LoadConfig()

	// init models
	models.Init(setting.IsProMode)

	// init social
	social.SetORM(models.ORM())
	setting.SocialAuth = social.NewSocial("/login/", auth.SocialAuther)
	setting.SocialAuth.ConnectSuccessURL = "/settings/profile"
	setting.SocialAuth.ConnectFailedURL = "/settings/profile"
	setting.SocialAuth.ConnectRegisterURL = "/register/connect"
	setting.SocialAuth.LoginURL = "/login"

	// init tango
	t := initTango(setting.IsProMode)

	// init routers
	routers.Init(t)

	// run
	setting.Log.Info("start WeGo", "v"+setting.APP_VER, setting.AppUrl)
	t.Run(setting.AppHost)
}

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
	. "github.com/qiniu/api/conf"
)

// We have to call a initialize function manully
// because we use `bee bale` to pack static resources
// and we cannot make sure that which init() execute first.
func initialize() {
	setting.LoadConfig()

	setting.SocialAuth = social.NewSocial("/login/", auth.SocialAuther)
	setting.SocialAuth.ConnectSuccessURL = "/settings/profile"
	setting.SocialAuth.ConnectFailedURL = "/settings/profile"
	setting.SocialAuth.ConnectRegisterURL = "/register/connect"
	setting.SocialAuth.LoginURL = "/login"

	//Qiniu
	ACCESS_KEY = setting.QiniuAccessKey
	SECRET_KEY = setting.QiniuSecurityKey
}

func initTango(isprod bool) *tango.Tango {
	middlewares.Init()

	tg := tango.NewWithLog(setting.Log)
	if isprod {
		tg.Mode = tango.Prod
	} else {
		tg.Mode = tango.Dev
	}

	if !isprod {
		tg.Use(debug.Debug(debug.Options{
			IgnorePrefix:     "/static",
			HideResponseBody: false,
			HideRequestBody:  false,
		}))
	}

	tg.Use(tango.ClassicHandlers...)

	tg.Use(
		tango.Static(tango.StaticOptions{
			RootPath: "./static",
			Prefix:   "static",
		}),
		tango.Static(tango.StaticOptions{
			RootPath: "./static_source",
			Prefix:   "static_source",
		}),
		session.New(time.Duration(setting.SessionCookieLifeTime)),
		middlewares.Renders,
		setting.Captcha,
	)
	if setting.EnableXSRF {
		tg.Use(xsrf.New(time.Duration(setting.SessionCookieLifeTime)))
	}
	tg.Use(flash.Flashes(), events.Events())
	return tg
}

func main() {
	// init config
	initialize()

	// init models
	models.Init(setting.IsProMode)

	// init tango
	t := initTango(setting.IsProMode)

	// initialize the routers
	routers.Init(t)

	// run
	setting.Log.Info("start WeGo", "v"+setting.APP_VER, setting.AppUrl)
	t.Run(setting.AppHost)
}

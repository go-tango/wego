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
	"html/template"

	"github.com/lunny/tango"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
	"github.com/go-tango/social-auth"

	"github.com/go-tango/wego/modules/models"
	"github.com/go-tango/wego/modules/utils"
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

	/*beego.SetLogFuncCall(true)*/
	setting.SocialAuth = social.NewSocial("/login/", auth.SocialAuther)
	setting.SocialAuth.ConnectSuccessURL = "/settings/profile"
	setting.SocialAuth.ConnectFailedURL = "/settings/profile"
	setting.SocialAuth.ConnectRegisterURL = "/register/connect"
	setting.SocialAuth.LoginURL = "/login"

	//Qiniu
	ACCESS_KEY = setting.QiniuAccessKey
	SECRET_KEY = setting.QiniuSecurityKey
}

func mergeFuncMap(funcs ...template.FuncMap) template.FuncMap{
	var ret = make(template.FuncMap)
	for _, fs := range funcs {
		for k, f := range fs {
			ret[k] = f
		}
	}
	return ret
}

func newTango() *tango.Tango {
	return tango.NewWithLog(
		setting.Log,
		tango.Logging(),
		tango.Recovery(true),
		tango.Compresses([]string{".js", ".css", ".html", ".htm"}),
		tango.Static("./static", "static", []string{"index.html", "index.htm"}),
		tango.Static("./static_source", "static_source", []string{"index.html", "index.htm"}),
		tango.Return(),
		tango.Responses(),
		tango.Requests(),
		tango.Param(),
		tango.Contexts(),
		session.New(time.Duration(setting.SessionCookieLifeTime)),
		renders.New(renders.Options{
			Funcs: mergeFuncMap(utils.FuncMap(), setting.Funcs),
		}),
	)
}

func main() {
	initialize()

	t := newTango()
	if setting.EnableXSRF {
		t.Use(xsrf.New(time.Duration(setting.SessionCookieLifeTime)))
	}
	t.Use(tango.Events())

	if setting.IsProMode {
		t.Mode = tango.Prod
	} else {
		t.Mode = tango.Dev
	}
	setting.Log.Info("start WeGo version", setting.APP_VER, setting.AppUrl)

	//initialize the routers
	routers.Init(t)

	// init models
	models.Init(setting.IsProMode)

	t.Run(setting.AppHost)
}

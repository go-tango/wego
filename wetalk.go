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
	"os"
	"fmt"
	"io"
	"time"
	"html/template"

	"github.com/lunny/tango"
	"github.com/lunny/log"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"

	"github.com/go-tango/wetalk/modules/models"
	"github.com/go-tango/wetalk/modules/utils"
	"github.com/go-tango/social-auth"
	"github.com/go-tango/wetalk/routers"
	"github.com/go-tango/wetalk/routers/auth"
	"github.com/go-tango/wetalk/setting"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/qiniu/api/conf"
)

type Prepare interface {
	Prepare()
}

type Finish interface {
	Finish()
}

// We have to call a initialize function manully
// because we use `bee bale` to pack static resources
// and we cannot make sure that which init() execute first.
func initialize() {
	setting.LoadConfig()

	//set logger
	f, err := os.Create("logs/wetalk.log")
	if err != nil {
		fmt.Println("create log file failed:", err)
		return
	}
	defer f.Close()

	w := io.MultiWriter(f, os.Stdout)
	log.SetOutput(w)

	if setting.IsProMode {
		log.SetOutputLevel(log.Linfo)
	} else {
		log.SetOutputLevel(log.Ldebug)
	}
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
	var logger = tango.NewLogger(os.Stdout)

	return tango.NewWithLog(
		logger,
		tango.NewLogging(logger),
		tango.NewRecovery(true),
		tango.NewCompress([]string{".js", ".css", ".html", ".htm"}),
		tango.NewStatic("./static", "static", []string{"index.html", "index.htm"}),
		tango.NewStatic("./static_source", "static_source", []string{"index.html", "index.htm"}),
		tango.HandlerFunc(tango.ReturnHandler),
		tango.HandlerFunc(tango.ResponseHandler),
		tango.HandlerFunc(tango.RequestHandler),
		tango.HandlerFunc(tango.ParamHandler),
		tango.HandlerFunc(tango.ContextHandler),
		tango.HandlerFunc(tango.EventHandler),
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
	t.Use(tango.HandlerFunc(func(ctx *tango.Context){
		if action := ctx.Action(); action != nil {
			if p, ok := action.(Prepare); ok {
				p.Prepare()
			}
		}
		ctx.Next()
	}))

	t.Use(tango.HandlerFunc(func(ctx *tango.Context){
		ctx.Next()

		if action := ctx.Action(); action != nil {
			if p, ok := action.(Finish); ok {
				p.Finish()
			}
		}
	}))

	if setting.IsProMode {
		t.Mode = tango.Prod
	} else {
		t.Mode = tango.Dev
	}
	log.Info(setting.APP_VER, setting.AppUrl)

	//initialize the routers
	routers.Init(t)

	// init models
	models.Init(setting.IsProMode)

	t.Run(setting.AppHost)
}

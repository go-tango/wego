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

package base

import (
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/i18n"
	"github.com/astaxie/beego/validation"
	"github.com/lunny/tango"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"

	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/auth"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
)

// baseRouter implemented global settings for all other routers.
type BaseRouter struct {
	//tango.Compress
	tango.Ctx
	session.Session
	xsrf.Checker
	renders.Renderer
	flash.Flash
	i18n.Locale

	User     models.User
	IsLogin  bool
	Data     renders.T
	TplNames string
}

// Before implemented Before method for baseRouter.
func (this *BaseRouter) Before() {
	this.Data = make(renders.T)

	if setting.EnforceRedirect {
		// if the host not matching app settings then redirect to AppUrl
		if this.Ctx.Req().Host != setting.AppHost {
			this.Redirect(setting.AppUrl)
			return
		}
	}

	// page start time
	this.Data["PageStartTime"] = time.Now()

	// check flash redirect, if match url then end, else for redirect return
	if match, redir := this.CheckFlashRedirect(this.Ctx.Req().RequestURI); redir {
		return
	} else if match {
		this.EndFlashRedirect()
	}

	switch {
	// save logined user if exist in session
	case auth.GetUserFromSession(&this.User, this.Session.Session):
		this.IsLogin = true
	// save logined user if exist in remember cookie
	case auth.LoginUserFromRememberCookie(&this.User, this.Ctx.Context, this.Session.Session):
		this.IsLogin = true
	}

	if this.IsLogin {
		this.IsLogin = true
		this.Data["User"] = &this.User
		this.Data["IsLogin"] = this.IsLogin

		// if user forbided then do logout
		if this.User.IsForbid {
			auth.LogoutUser(this.Context, this.Session.Session)
			this.FlashRedirect("/login", 302, "UserForbid")
			return
		}
	}

	// Setting properties.
	this.Data["AppName"] = setting.AppName
	this.Data["AppVer"] = setting.AppVer
	this.Data["AppUrl"] = setting.AppUrl
	this.Data["AppLogo"] = setting.AppLogo
	this.Data["AvatarURL"] = setting.AvatarURL
	this.Data["IsProMode"] = setting.IsProMode
	this.Data["SearchEnabled"] = setting.SearchEnabled
	this.Data["Flush"] = this.Flash.Data()

	// Redirect to make URL clean.
	if this.setLang() {
		i := strings.Index(this.Ctx.Req().RequestURI, "?")
		this.Redirect(this.Ctx.Req().RequestURI[:i])
		return
	}

	// pass xsrf helper to template context
	this.Data["xsrf_token"] = this.XsrfValue
	this.Data["xsrf_html"] = this.XsrfFormHtml()

	// read unread notifications
	if this.IsLogin {
		this.Data["UnreadNotificationCount"] = models.GetUnreadNotificationCount(this.User.Id)
	}

	// if method is GET then auto create a form once token
	if this.Ctx.Req().Method == "GET" {
		this.FormOnceCreate()
	}
}

// on router finished
func (this *BaseRouter) After() {
	if !this.Ctx.Written() && this.TplNames != "" {
		err := this.Render(this.TplNames, this.Data)
		if err != nil {
			this.Result = err
		}
	}
}

func (this *BaseRouter) LoginUser(user *models.User, remember bool) string {
	loginRedirect := strings.TrimSpace(auth.GetCookie(this.Req(), "login_to"))
	if utils.IsMatchHost(loginRedirect) == false {
		loginRedirect = "/"
	} else {
		auth.SetCookie(this, "login_to", "", -1, "/")
	}

	// login user
	auth.LoginUser(user, this.Context, this.Session.Session, remember)

	this.setLangCookie(i18n.GetLangByIndex(user.Lang))

	return loginRedirect
}

// check if user not active then redirect
func (this *BaseRouter) CheckActiveRedirect(args ...interface{}) bool {
	var redirect_to string
	code := 302
	needActive := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needActive = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			code = v
		}
	}

	if needActive {
		// check login
		if this.CheckLoginRedirect() {
			return true
		}

		// redirect to active page
		if !this.User.IsActive {
			this.FlashRedirect("/settings/profile", code, "NeedActive")
			return true
		}
	} else {
		// no need active
		if this.User.IsActive {
			if redirect_to == "" {
				redirect_to = "/"
			}
			this.Redirect(redirect_to, code)
			return true
		}
	}
	return false
}

// check if not login then redirect
func (this *BaseRouter) CheckLoginRedirect(args ...interface{}) bool {
	var redirect_to string
	code := 302
	needLogin := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needLogin = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			// custom redirect url
			code = v
		}
	}

	// if need login then redirect
	if needLogin && !this.IsLogin {
		if len(redirect_to) == 0 {
			req := this.Ctx.Req()
			scheme := "http"
			if req.TLS != nil {
				scheme += "s"
			}
			redirect_to = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)
		}
		redirect_to = "/login?to=" + url.QueryEscape(redirect_to)
		this.Redirect(redirect_to, code)
		return true
	}

	// if not need login then redirect
	if !needLogin && this.IsLogin {
		if len(redirect_to) == 0 {
			redirect_to = "/"
		}
		this.Redirect(redirect_to, code)
		return true
	}
	return false
}

// check flash redirect, ensure browser redirect to uri and display flash message.
func (this *BaseRouter) CheckFlashRedirect(value string) (match bool, redirect bool) {
	v := this.Session.Get("on_redirect")
	if params, ok := v.([]interface{}); ok {
		if len(params) != 5 {
			this.EndFlashRedirect()
			goto end
		}
		uri := utils.ToStr(params[0])
		code := 302
		if c, ok := params[1].(int); ok {
			if c/100 == 3 {
				code = c
			}
		}
		flag := utils.ToStr(params[2])
		flagVal := utils.ToStr(params[3])
		times := 0
		if v, ok := params[4].(int); ok {
			times = v
		}

		times += 1
		if times > 3 {
			// if max retry times reached then end
			this.EndFlashRedirect()
			goto end
		}

		// match uri or flash flag
		if uri == value || flag == value {
			match = true
		} else {
			// if no match then continue redirect
			this.FlashRedirect(uri, code, flag, flagVal, times)
			redirect = true
		}
	}
end:
	return match, redirect
}

// set flash redirect
func (this *BaseRouter) FlashRedirect(uri string, code int, flag string, args ...interface{}) {
	flagVal := "true"
	times := 0
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			flagVal = v
		case int:
			times = v
		}
	}

	if len(uri) == 0 || uri[0] != '/' {
		panic("flash reirect only support same host redirect")
	}

	params := []interface{}{uri, code, flag, flagVal, times}
	this.Session.Set("on_redirect", params)

	this.FlashWrite(flag, flagVal)
	this.Flash.Redirect(uri)
}

func (this *BaseRouter) FlashWrite(key, value string) {
	this.Flash.Set(key, value)
	this.Data["flash"] = this.Flash.Data()
}

// clear flash redirect
func (this *BaseRouter) EndFlashRedirect() {
	this.Session.Del("on_redirect")
}

// check form once, void re-submit
func (this *BaseRouter) FormOnceNotMatch() bool {
	notMatch := false
	recreat := false

	// get token from request param / header
	var value string
	if vus := this.Req().FormValue("_once"); len(vus) > 0 {
		value = vus
	} else {
		value = this.Ctx.Header().Get("X-Form-Once")
	}

	// exist in session
	if v, ok := this.Session.Get("form_once").(string); ok && v != "" {
		// not match
		if value != v {
			notMatch = true
		} else {
			// if matched then re-creat once
			recreat = true
		}
	}

	this.FormOnceCreate(recreat)
	return notMatch
}

func (this *BaseRouter) GetSession(key string) interface{} {
	return this.Session.Get(key)
}

// create form once html
func (this *BaseRouter) FormOnceCreate(args ...bool) {
	var value string
	var creat bool
	creat = len(args) > 0 && args[0]
	if !creat {
		if v, ok := this.Session.Get("form_once").(string); ok && v != "" {
			value = v
		} else {
			creat = true
		}
	}
	if creat {
		value = utils.GetRandomString(10)
		this.Session.Set("form_once", value)
	}
	this.Data["once_token"] = value
	this.Data["once_html"] = template.HTML(`<input type="hidden" name="_once" value="` + value + `">`)
}

func (this *BaseRouter) validForm(form interface{}, names ...string) (bool, map[string]*validation.ValidationError) {
	// parse request params to form ptr struct
	utils.ParseForm(form, this.Ctx.Req().Form)

	// Put data back in case users input invalid data for any section.
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	this.Data[name] = form

	errName := name + "Error"

	// check form once
	if this.FormOnceNotMatch() {
		return false, nil
	}

	// Verify basic input.
	valid := validation.Validation{}
	if ok, _ := valid.Valid(form); !ok {
		errs := valid.ErrorMap()
		this.Data[errName] = &valid
		return false, errs
	}
	return true, nil
}

// valid form and put errors to tempalte context
func (this *BaseRouter) ValidForm(form interface{}, names ...string) bool {
	valid, _ := this.validForm(form, names...)
	return valid
}

// valid form and put errors to tempalte context
func (this *BaseRouter) ValidFormSets(form interface{}, names ...string) bool {
	valid, errs := this.validForm(form, names...)
	this.setFormSets(form, errs, names...)
	return valid
}

func (this *BaseRouter) SetFormSets(form interface{}, names ...string) *utils.FormSets {
	return this.setFormSets(form, nil, names...)
}

func (this *BaseRouter) setFormSets(form interface{}, errs map[string]*validation.ValidationError, names ...string) *utils.FormSets {
	formSets := utils.NewFormSets(form, errs, this.Locale)
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	name += "Sets"
	this.Data[name] = formSets

	return formSets
}

// add valid error to FormError
func (this *BaseRouter) SetFormError(form interface{}, fieldName, errMsg string, names ...string) {
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	errName := name + "Error"
	setsName := name + "Sets"

	if valid, ok := this.Data[errName].(*validation.Validation); ok {
		valid.SetError(fieldName, this.Tr(errMsg))
	}

	if fSets, ok := this.Data[setsName].(*utils.FormSets); ok {
		fSets.SetError(fieldName, errMsg)
	}
}

func (this *BaseRouter) IsAjax() bool {
	return this.Req().Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (this *BaseRouter) SetPaginator(per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(this.Req(), per, nums)
	this.Data["paginator"] = p
	return p
}

func (this *BaseRouter) JsStorage(action, key string, values ...string) {
	value := action + ":::" + key
	if len(values) > 0 {
		value += ":::" + values[0]
	}
	auth.SetCookie(this, "JsStorage", value, 1<<31-1, "/", nil, nil, false)
}

func (this *BaseRouter) setLangCookie(lang string) {
	auth.SetCookie(this, "lang", lang, 60*60*24*365, "/", nil, nil, false)
}

func (this *BaseRouter) GetString(k string) string {
	return this.Req().FormValue(k)
}

func (this *BaseRouter) GetInt(k string) (int64, error) {
	return strconv.ParseInt(this.Req().FormValue(k), 10, 64)
}

// setLang sets site language version.
func (this *BaseRouter) setLang() bool {
	isNeedRedir := false
	hasCookie := false

	// get all lang names from i18n
	langs := setting.Langs

	// 1. Check URL arguments.
	lang := this.Req().FormValue("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = auth.GetCookie(this.Req(), "lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. check if isLogin then use user setting
	if len(lang) == 0 && this.IsLogin {
		lang = i18n.GetLangByIndex(this.User.Lang)
	}

	// 4. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := this.Header().Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. DefaultLang language is Chinese.
	if len(lang) == 0 {
		//lang = "en-US"
		lang = "zh-CN"
		isNeedRedir = false
	}

	// Save language information in cookies.
	if !hasCookie {
		this.setLangCookie(lang)
	}

	// Set language properties.
	this.Data["Lang"] = lang
	this.Data["Langs"] = langs

	this.Lang = lang

	return isNeedRedir
}

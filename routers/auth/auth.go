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

package auth

import (
	"strings"

	"github.com/lunny/log"

	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/auth"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/routers/base"
	"github.com/go-tango/wego/setting"
)

// LoginRouter serves login page.
type Login struct {
	base.BaseRouter
}

// Get implemented login page.
func (this *Login) Get() {
	this.Data["IsLoginPage"] = true

	loginRedirect := strings.TrimSpace(this.GetString("to"))
	if loginRedirect == "" {
		loginRedirect = this.Ctx.Header().Get("Referer")
	}
	if utils.IsMatchHost(loginRedirect) == false {
		loginRedirect = "/"
	}

	// no need login
	if this.CheckLoginRedirect(false, loginRedirect) {
		return
	}

	if len(loginRedirect) > 0 {
		auth.SetCookie(this, "login_to", loginRedirect, 0, "/")
	}

	form := auth.LoginForm{}
	this.SetFormSets(&form)

	this.Render("auth/login.html", this.Data)
}

// Login implemented user login.
func (this *Login) Post() {
	this.Data["IsLoginPage"] = true

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	var user models.User
	var key string
	ajaxErrMsg := "auth.login_error_ajax"

	form := auth.LoginForm{}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		if this.IsAjax() {
			goto ajaxError
		}
		return
	}

	key = "auth.login." + form.UserName + utils.IP(this.Req())
	if times, ok := utils.TimesReachedTest(key, setting.LoginMaxRetries); ok {
		if this.IsAjax() {
			ajaxErrMsg = "auth.login_error_times_reached"
			goto ajaxError
		}
		this.Data["ErrorReached"] = true

	} else if auth.VerifyUser(&user, form.UserName, form.Password) {
		loginRedirect := this.LoginUser(&user, form.Remember)

		if this.IsAjax() {
			this.Data["json"] = map[string]interface{}{
				"success":  true,
				"message":  this.Tr("auth.login_success_ajax"),
				"redirect": loginRedirect,
			}
			this.ServeJson(this.Data)
			return
		}

		this.Redirect(loginRedirect, 302)
		return
	} else {
		utils.TimesReachedSet(key, times, setting.LoginFailedBlocks)
		if this.IsAjax() {
			goto ajaxError
		}
	}
	this.Data["Error"] = true
	this.Render("auth/login.html", this.Data)
	return

ajaxError:
	this.Data["json"] = map[string]interface{}{
		"success": false,
		"message": this.Tr(ajaxErrMsg),
		"once":    this.Data["once_token"],
	}
	this.ServeJson(this.Data)
}

type Logout struct {
	base.BaseRouter
}

// Logout implemented user logout page.
func (this *Logout) Get() {
	auth.LogoutUser(this.Context, this.Session.Session)

	// write flash message
	//this.FlashWrite("HasLogout", "true")

	this.Redirect("/login")
}

// RegisterRouter serves register page.
type Register struct {
	base.BaseRouter
}

// Get implemented Get method for RegisterRouter.
func (this *Register) Get() {
	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	this.Data["IsRegisterPage"] = true

	form := auth.RegisterForm{Locale: this.Locale}
	this.SetFormSets(&form)
	this.Render("auth/register.html", this.Data)
}

// Register implemented Post method for RegisterRouter.
func (this *Register) Post() {
	this.Data["IsRegisterPage"] = true

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	form := auth.RegisterForm{Locale: this.Locale}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		return
	}

	// Create new user.
	user := new(models.User)

	if err := auth.RegisterUser(user, form.UserName, form.Email, form.Password, this.Locale); err == nil {
		auth.SendRegisterMail(this.Locale, user)

		loginRedirect := this.LoginUser(user, false)
		if loginRedirect == "/" {
			this.FlashRedirect("/settings/profile", 302, "RegSuccess")
		} else {
			this.Redirect(loginRedirect)
			return
		}

		this.Render("auth/register.html", this.Data)
	} else {
		log.Error("Register: Failed ", err)
	}
}

type RegisterActive struct {
	base.BaseRouter
}

// Active implemented check Email actice code.
func (this *RegisterActive) Get() {
	// no need active
	if this.CheckActiveRedirect(false) {
		return
	}

	code := this.Params().Get(":code")

	var user models.User

	if auth.VerifyUserActiveCode(&user, code) {
		user.IsActive = true
		user.Rands = models.GetUserSalt()
		if err := models.UpdateById(user.Id, user, models.Obj2Table([]string{"IsActive", "Rands", "Updated"})...); err != nil {
			log.Error("Active: user Update ", err)
		}
		if this.IsLogin {
			this.User = user
		}

		this.Redirect("/active/success", 302)

	} else {
		this.Data["Success"] = false
	}

	this.Render("auth/active.html", this.Data)
}

type RegisterSuccess struct {
	base.BaseRouter
}

// ActiveSuccess implemented success page when email active code verified.
func (this *RegisterSuccess) Get() {
	this.Data["Success"] = true
	this.Render("auth/active.html", this.Data)
}

// ForgotRouter serves login page.
type ForgotRouter struct {
	base.BaseRouter
}

// Get implemented Get method for ForgotRouter.
func (this *ForgotRouter) Get() {
	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	form := auth.ForgotForm{}
	this.SetFormSets(&form)
	this.Render("auth/forgot.html", this.Data)
}

// Get implemented Post method for ForgotRouter.
func (this *ForgotRouter) Post() {
	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	var user models.User
	form := auth.ForgotForm{User: &user}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		return
	}

	// send reset password email
	auth.SendResetPwdMail(this.Locale, &user)

	this.FlashRedirect("/forgot", 302, "SuccessSend")

	this.Render("auth/forgot.html", this.Data)
}

// ForgotRouter serves login page.
type ResetRouter struct {
	base.BaseRouter
}

// Reset implemented user password reset.
func (this *ResetRouter) Get() {
	code := this.GetString(":code")
	this.Data["Code"] = code

	var user models.User

	if auth.VerifyUserResetPwdCode(&user, code) {
		this.Data["Success"] = true
		form := auth.ResetPwdForm{}
		this.SetFormSets(&form)
	} else {
		this.Data["Success"] = false
	}
	this.Render("auth/reset.html", this.Data)
}

// Reset implemented user password reset.
func (this *ResetRouter) Post() {
	code := this.GetString(":code")
	this.Data["Code"] = code

	var user models.User

	if auth.VerifyUserResetPwdCode(&user, code) {
		this.Data["Success"] = true

		form := auth.ResetPwdForm{}
		if this.ValidFormSets(&form) == false {
			return
		}

		user.IsActive = true
		user.Rands = models.GetUserSalt()

		if err := auth.SaveNewPassword(&user, form.Password); err != nil {
			log.Error("ResetPost Save New Password: ", err)
		}

		if this.IsLogin {
			auth.LogoutUser(this.Context, this.Session.Session)
			return
		}

		this.FlashRedirect("/login", 302, "ResetSuccess")

	} else {
		this.Data["Success"] = false
	}

	this.Render("auth/reset.html", this.Data)
}

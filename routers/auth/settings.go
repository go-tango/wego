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
	"github.com/lunny/log"

	"github.com/go-tango/wego/modules/auth"
	"github.com/go-tango/wego/routers/base"
	"github.com/go-tango/wego/setting"
)

// SettingsRouter serves user settings.
type ProfileRouter struct {
	base.BaseRouter
}

// Profile implemented user profile settings page.
func (this *ProfileRouter) Get() error {
	this.Data["IsUserSettingPage"] = true

	// need login
	if this.CheckLoginRedirect() {
		return nil
	}

	form := auth.ProfileForm{Locale: this.Locale}

	form.SetFromUser(&this.User)
	this.SetFormSets(&form)
	return this.Render("settings/profile.html", this.Data)
}

// ProfileSave implemented save user profile.
func (this *ProfileRouter) Post() {
	this.Data["IsUserSettingPage"] = true
	if this.CheckLoginRedirect() {
		return
	}

	action := this.GetString("action")
	if this.IsAjax() {
		switch action {
		case "send-verify-email":
			if this.User.IsActive {
				this.Data["json"] = false
			} else {
				auth.SendActiveMail(this.Locale, &this.User)
				this.Data["json"] = true
			}

			this.ServeJson(this.Data)
			return
		}
		return
	}

	profileForm := auth.ProfileForm{Locale: this.Locale}
	profileForm.SetFromUser(&this.User)

	this.Data["Form"] = profileForm

	if this.ValidFormSets(&profileForm) {
		if err := profileForm.SaveUserProfile(&this.User); err != nil {
			log.Error("ProfileSave: save-profile", err)
		}
		this.FlashRedirect("/settings/profile", 302, "ProfileSave")
		return
	}

	this.Render("settings/profile.html", this.Data)
}

// SettingsRouter serves user settings.
type PasswordRouter struct {
	base.BaseRouter
}

func (this *PasswordRouter) Get() error {
	this.Data["IsUserSettingPage"] = true

	//need login
	if this.CheckLoginRedirect() {
		return nil
	}

	formPwd := auth.PasswordForm{}
	this.SetFormSets(&formPwd)
	return this.Render("settings/change_password.html", this.Data)
}

func (this *PasswordRouter) Post() {
	this.Data["IsUserSettingPage"] = true
	if this.CheckLoginRedirect() {
		return
	}

	pwdForm := auth.PasswordForm{User: &this.User}

	this.Data["Form"] = pwdForm

	if this.ValidFormSets(&pwdForm) {
		// verify success and save new password
		if err := auth.SaveNewPassword(&this.User, pwdForm.Password); err == nil {
			this.FlashRedirect("/settings/change/password", 302, "PasswordSave")
			return
		} else {
			log.Error("ProfileSave: change-password", err)
		}
	}

	this.Render("settings/change_password.html", this.Data)
}

// SettingsRouter serves user settings.
type AvatarRouter struct {
	base.BaseRouter
}

func (this *AvatarRouter) Get() error {
	this.Data["IsUserSettingPage"] = true
	//need login
	if this.CheckLoginRedirect() {
		return nil
	}

	form := auth.UserAvatarForm{}
	form.SetFromUser(&this.User)
	this.SetFormSets(&form)
	return this.Render("settings/user_avatar.html", this.Data)
}

func (this *AvatarRouter) Post() {
	this.Data["IsUserSettingPage"] = true
	//need login
	if this.CheckLoginRedirect() {
		return
	}
	avatarType, _ := this.GetInt("AvatarType")
	form := auth.UserAvatarForm{AvatarType: int(avatarType)}
	this.Data["Form"] = form

	if this.ValidFormSets(&form) {
		if err := auth.SaveAvatarType(&this.User, int(avatarType)); err == nil {
			this.FlashRedirect("/settings/avatar", 302, "AvatarSettingSave")
			return
		} else {
			log.Error("ProfileSave: avatar-setting", err)
		}
	}
	this.Render("settings/user_avatar.html", this.Data)
}

type AvatarUploadRouter struct {
	base.BaseRouter
}

func (this *AvatarUploadRouter) Post() {
	this.Data["IsUserSettingPage"] = true
	//need login and active
	if this.CheckLoginRedirect() {
		return
	}

	// get file object
	file, handler, err := this.Ctx.Req().FormFile("avatar")
	if err != nil {
		return
	}
	defer file.Close()
	mime := handler.Header.Get("Content-Type")
	if err := auth.UploadUserAvatarToQiniu(file, handler.Filename, mime, setting.QiniuAvatarBucket, &this.User); err != nil {
		return
	}

	userAvatarForm := auth.UserAvatarForm{}
	userAvatarForm.SetFromUser(&this.User)
	this.SetFormSets(&userAvatarForm)
	this.FlashRedirect("/settings/avatar", 302, "AvatarUploadSuccess")
	this.Render("settings/user_avatar.html", this.Data)
}

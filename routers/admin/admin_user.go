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

package admin

import (
	"fmt"

	"github.com/lunny/log"
	"github.com/astaxie/beego/orm"

	"github.com/go-tango/wego/modules/auth"
	"github.com/go-tango/wego/modules/models"
	"github.com/go-tango/wego/modules/utils"
)

type UserAdminRouter struct {
	ModelAdminRouter
	object models.User
}

func (this *UserAdminRouter) Before() {
	this.Params().Set(":model", "user")
	this.ModelAdminRouter.Before()
}

func (this *UserAdminRouter) Object() interface{} {
	return &this.object
}

func (this *UserAdminRouter) ObjectQs() orm.QuerySeter {
	return models.Users().RelatedSel()
}

type UserAdminList struct {
	UserAdminRouter
}

// view for list model data
func (this *UserAdminList) Get() {
	var q = this.GetString("q")
	var users []models.User
	var qs orm.QuerySeter
	if q != "" {
		cond := orm.NewCondition()
		cond = cond.Or("Email", q)
		cond = cond.Or("UserName", q)
		qs = models.Users().SetCond(cond)
	} else {
		qs = models.Users()
	}
	this.Data["q"] = q
	if err := this.SetObjects(qs, &users); err != nil {
		this.Data["Error"] = err
		log.Error(err)
	}
}

type UserAdminNew struct {
	UserAdminRouter
}

// view for create object
func (this *UserAdminNew) Get() {
	form := auth.UserAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *UserAdminNew) Post() {
	form := auth.UserAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var user models.User
	form.SetToUser(&user)
	if err := user.Insert(); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/user/%d", user.Id), 302, "CreateSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

type UserAdminEdit struct {
	UserAdminRouter
}

// view for edit object
func (this *UserAdminEdit) Get() {
	form := auth.UserAdminForm{}
	form.SetFromUser(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *UserAdminEdit) Post() {
	form := auth.UserAdminForm{Id: this.object.Id}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/user/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToUser(&this.object)
		if err := this.object.Update(changes...); err == nil {
			this.FlashRedirect(url, 302, "UpdateSuccess")
			return
		} else {
			log.Error(err)
			this.Data["Error"] = err
		}
	} else {
		this.Redirect(url, 302)
	}
}

type UserAdminDelete struct {
	UserAdminRouter
}

// view for delete object
func (this *UserAdminDelete) Post() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := this.object.Delete(); err == nil {
		this.FlashRedirect("/admin/user", 302, "DeleteSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

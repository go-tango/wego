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

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/go-tango/wetalk/modules/models"
	"github.com/go-tango/wetalk/modules/page"
	"github.com/go-tango/wetalk/modules/utils"
)

type PageAdminRouter struct {
	ModelAdminRouter
	object models.Page
}

func (this *PageAdminRouter) Object() interface{} {
	return &this.object
}

func (this *PageAdminRouter) ObjectQs() orm.QuerySeter {
	return models.Pages().RelatedSel()
}

// view for list model data
func (this *PageAdminRouter) List() {
	var pages []models.Page
	qs := models.Pages().RelatedSel()
	if err := this.SetObjects(qs, &pages); err != nil {
		this.Data["Error"] = err
		beego.Error(err)
	}
}

// view for create object
func (this *PageAdminRouter) Create() {
	form := page.PageAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *PageAdminRouter) Save() {
	form := page.PageAdminForm{Create: true}
	if !this.ValidFormSets(&form) {
		return
	}

	var a models.Page
	form.SetToPage(&a)
	if err := a.Insert(); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/page/%d", a.Id), 302, "CreateSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}

// view for edit object
func (this *PageAdminRouter) Edit() {
	form := page.PageAdminForm{}
	form.SetFromPage(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *PageAdminRouter) Update() {
	form := page.PageAdminForm{}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/page/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToPage(&this.object)
		if err := this.object.Update(changes...); err == nil {
			this.FlashRedirect(url, 302, "UpdateSuccess")
			return
		} else {
			beego.Error(err)
			this.Data["Error"] = err
		}
	} else {
		this.Redirect(url, 302)
	}
}

// view for confirm delete object
func (this *PageAdminRouter) Confirm() {
}

// view for delete object
func (this *PageAdminRouter) Delete() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := this.object.Delete(); err == nil {
		this.FlashRedirect("/admin/page", 302, "DeleteSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}

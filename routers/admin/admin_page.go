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

	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/page"
	"github.com/go-tango/wego/modules/utils"
)

type PageAdminRouter struct {
	ModelAdminRouter
	object models.Page
}

func (this *PageAdminRouter) Before() {
	this.Params().Set(":model", "page")
	this.ModelAdminRouter.Before()
}

func (this *PageAdminRouter) Object() interface{} {
	return &this.object
}

type PageAdminList struct {
	PageAdminRouter
}

// view for list model data
func (this *PageAdminList) Get() {
	var pages []models.Page
	sess := models.Orm().NewSession()
	defer sess.Close()
	if err := this.SetObjects(sess, &pages); err != nil {
		this.Data["Error"] = err
		log.Error(err)
	}
}

type PageAdminNew struct {
	PageAdminRouter
}

// view for create object
func (this *PageAdminNew) Get() {
	form := page.PageAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *PageAdminNew) Post() {
	form := page.PageAdminForm{Create: true}
	if !this.ValidFormSets(&form) {
		return
	}

	var a models.Page
	form.SetToPage(&a)
	if err := models.Insert(&a); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/page/%d", a.Id), 302, "CreateSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

type PageAdminEdit struct {
	PageAdminRouter
}

// view for edit object
func (this *PageAdminEdit) Get() {
	form := page.PageAdminForm{}
	form.SetFromPage(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *PageAdminEdit) Post() {
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
		if err := models.UpdateById(this.object.Id, this.object, models.Obj2Table(changes)...); err == nil {
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

type PageAdminDelete struct {
	PageAdminRouter
}

// view for delete object
func (this *PageAdminDelete) Post() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := models.DeleteById(this.object.Id, this.object); err == nil {
		this.FlashRedirect("/admin/page", 302, "DeleteSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

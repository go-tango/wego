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
	"github.com/go-tango/wego/modules/post"
	"github.com/go-tango/wego/modules/utils"
)

type CategoryAdminRouter struct {
	ModelAdminRouter
	object models.Category
}

func (this *CategoryAdminRouter) Before() {
	this.Params().Set(":model", "category")
	this.ModelAdminRouter.Before()
}

func (this *CategoryAdminRouter) Object() interface{} {
	return &this.object
}

type CategoryAdminList struct {
	CategoryAdminRouter
}

// view for list model data
func (this *CategoryAdminList) Get() {
	var cats []models.Category
	sess := models.ORM().NewSession()
	defer sess.Close()
	if err := this.SetObjects(sess, &cats); err != nil {
		this.Data["Error"] = err
		log.Error(err)
	}
}

type CategoryAdminNew struct {
	CategoryAdminRouter
}

// view for create object
func (this *CategoryAdminNew) Get() {
	form := post.CategoryAdminForm{Create: true}
	this.SetFormSets(&form)
}

type CategoryAdminEdit struct {
	CategoryAdminRouter
}

// view for edit object
func (this *CategoryAdminEdit) Get() {
	form := post.CategoryAdminForm{}
	form.SetFromCategory(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *CategoryAdminEdit) Post() {
	form := post.CategoryAdminForm{Id: int(this.object.Id)}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/category/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToCategory(&this.object)
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

type CategoryAdminDelete struct {
	CategoryAdminRouter
}

// view for delete object
func (this *CategoryAdminDelete) Post() {
	if this.FormOnceNotMatch() {
		return
	}
	// check whether there are topics under the category
	cnt, _ := models.CountTopicsByCategoryId(this.object.Id)
	if cnt > 0 {
		this.FlashRedirect("/admin/category", 302, "DeleteNotAllowed")
		return
	} else {
		// delete object
		if err := models.DeleteById(this.object.Id, this.object); err == nil {
			this.FlashRedirect("/admin/category", 302, "DeleteSuccess")
			return
		} else {
			log.Error(err)
			this.Data["Error"] = err
		}
	}
}

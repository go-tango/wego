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

	"github.com/astaxie/beego/orm"
	"github.com/lunny/log"

	"github.com/go-tango/wego/modules/models"
	"github.com/go-tango/wego/modules/post"
	"github.com/go-tango/wego/modules/utils"
)

type PostAdminRouter struct {
	ModelAdminRouter
	object models.Post
}

func (this *PostAdminRouter) Before() {
	this.Params().Set(":model", "post")
	this.ModelAdminRouter.Before()
}

func (this *PostAdminRouter) Object() interface{} {
	return &this.object
}

func (this *PostAdminRouter) ObjectQs() orm.QuerySeter {
	return models.Posts().RelatedSel()
}

func (this *PostAdminRouter) GetForm(create bool) post.PostAdminForm {
	form := post.PostAdminForm{Create: create}
	post.ListTopics(&form.Topics)
	return form
}

type PostAdminList struct {
	PostAdminRouter
}

// view for list model data
func (this *PostAdminList) Get() {
	var posts []models.Post
	qs := models.Posts().RelatedSel()
	if err := this.SetObjects(qs, &posts); err != nil {
		this.Data["Error"] = err
		log.Error(err)
	}
}

type PostAdminNew struct {
	PostAdminRouter
}

// view for create object
func (this *PostAdminNew) Get() {
	form := this.GetForm(true)
	this.SetFormSets(&form)
}

// view for new object save
func (this *PostAdminNew) Post() {
	form := this.GetForm(true)
	if !this.ValidFormSets(&form) {
		return
	}

	var post models.Post
	form.SetToPost(&post)
	if err := post.Insert(); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/post/%d", post.Id), 302, "CreateSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

type PostAdminEdit struct {
	PostAdminRouter
}

// view for edit object
func (this *PostAdminEdit) Get() {
	form := this.GetForm(false)
	form.SetFromPost(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *PostAdminEdit) Post() {
	form := this.GetForm(false)
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/post/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		//fix the bug of category not updated
		changes = append(changes, "Category")
		form.SetToPost(&this.object)
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

type PostAdminDelete struct {
	PostAdminRouter
}

// view for delete object
func (this *PostAdminDelete) Get() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := this.object.Delete(); err == nil {
		this.FlashRedirect("/admin/post", 302, "DeleteSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

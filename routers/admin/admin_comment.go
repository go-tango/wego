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

type CommentAdminRouter struct {
	ModelAdminRouter
	object models.Comment
}

func (this *CommentAdminRouter) Before() {
	this.Params().Set(":model", "comment")
	this.ModelAdminRouter.Before()
}

func (this *CommentAdminRouter) Object() interface{} {
	return &this.object
}

type CommentAdminList struct {
	CommentAdminRouter
}

// view for list model data
func (this *CommentAdminList) Get() {
	var comments []models.Comment
	sess := models.ORM().NewSession()
	defer sess.Close()
	if err := this.SetObjects(sess, &comments); err != nil {
		this.Data["Error"] = err
		log.Error(err)
	}
}

type CommentAdminNew struct {
	CommentAdminRouter
}

// view for create object
func (this *CommentAdminNew) Get() {
	form := post.CommentAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *CommentAdminNew) Post() {
	form := post.CommentAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var comment models.Comment
	form.SetToComment(&comment)
	if err := models.Insert(&comment); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/comment/%d", comment.Id), 302, "CreateSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

type CommentAdminEdit struct {
	CommentAdminRouter
}

// view for edit object
func (this *CommentAdminEdit) Get() {
	form := post.CommentAdminForm{}
	form.SetFromComment(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *CommentAdminEdit) Post() {
	form := post.CommentAdminForm{}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/comment/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToComment(&this.object)
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

type CommentAdminDelete struct {
	CommentAdminRouter
}

// view for delete object
func (this *CommentAdminDelete) Post() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := models.DeleteById(this.object.Id, this.object); err == nil {
		this.FlashRedirect("/admin/comment", 302, "DeleteSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

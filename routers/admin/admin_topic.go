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

type TopicAdminRouter struct {
	ModelAdminRouter
	object models.Topic
}

func (this *TopicAdminRouter) Before() {
	this.Params().Set(":model", "topic")
	this.ModelAdminRouter.Before()
}

func (this *TopicAdminRouter) Object() interface{} {
	return &this.object
}

type TopicAdminList struct {
	TopicAdminRouter
}

// view for list model data
func (this *TopicAdminList) Get() {
	var topics []models.Topic
	sess := models.Orm().Desc("category_id")
	if err := this.SetObjects(sess, &topics); err != nil {
		this.Data["Error"] = err
		log.Error(err)
	}
}

type TopicAdminNew struct {
	TopicAdminRouter
}

// view for create object
func (this *TopicAdminNew) Get() {
	form := post.TopicAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *TopicAdminNew) Post() {
	form := post.TopicAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var topic models.Topic
	form.SetToTopic(&topic)
	if err := models.Insert(&topic); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/topic/%d", topic.Id), 302, "CreateSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

type TopicAdminEdit struct {
	TopicAdminRouter
}

// view for edit object
func (this *TopicAdminEdit) Get() {
	form := post.TopicAdminForm{}
	form.SetFromTopic(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *TopicAdminEdit) Post() {
	form := post.TopicAdminForm{Id: int(this.object.Id)}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/topic/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToTopic(&this.object)
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

type TopicAdminDelete struct {
	TopicAdminRouter
}

// view for delete object
func (this *TopicAdminDelete) Post() {
	if this.FormOnceNotMatch() {
		return
	}
	//check whether there are posts under this topic
	cnt, _ := models.Count(&models.Post{TopicId: this.object.Id})
	if cnt > 0 {
		this.FlashRedirect("/admin/topic", 302, "DeleteNotAllowed")
		return
	} else {
		// delete object
		if err := models.DeleteById(this.object.Id, this.object); err == nil {
			this.FlashRedirect("/admin/topic", 302, "DeleteSuccess")
			return
		} else {
			log.Error(err)
			this.Data["Error"] = err
		}
	}
}

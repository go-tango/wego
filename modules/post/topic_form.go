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

package post

import (
	"github.com/astaxie/beego/validation"

	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/utils"
)

type TopicAdminForm struct {
	Create    bool   `form:"-"`
	Id        int    `form:"-"`
	Name      string `valid:"Required;MaxSize(30)"`
	Intro     string `form:"type(textarea)" valid:"Required"`
	Slug      string `valid:"Required;MaxSize(100)"`
	Followers int    `form:"-"`
	Order     int    ``
	ImageLink string `valid:"MaxSize(200)"`
	Category  int    `form:"type(select);attr(rel,select2)" valid:""`
}

func (form *TopicAdminForm) CategorySelectData() [][]string {
	var cats []models.Category
	models.FindCategories(&cats)
	data := make([][]string, 0, len(cats))
	for _, cat := range cats {
		data = append(data, []string{cat.Name, utils.ToStr(cat.Id)})
	}
	return data
}

func (form *TopicAdminForm) Labels() map[string]string {
	return map[string]string{
		"Name":      "model.topic_name",
		"Intro":     "model.topic_intro",
		"Slug":      "model.topic_slug",
		"Order":     "model.topic_order",
		"ImageLink": "model.topic_image_link",
		"Category":  "model.category",
	}
}

func (form *TopicAdminForm) Valid(v *validation.Validation) {
	if models.IsExist(&models.Topic{Name: form.Name, Id: int64(form.Id)}) {
		v.SetError("Name", "admin.field_need_unique")
	}

	if models.IsExist(&models.Topic{Slug: form.Slug, Id: int64(form.Id)}) {
		v.SetError("Slug", "admin.field_need_unique")
	}
}

func (form *TopicAdminForm) SetFromTopic(topic *models.Topic) {
	utils.SetFormValues(topic, form)
	form.Category = int(topic.CategoryId)
}

func (form *TopicAdminForm) SetToTopic(topic *models.Topic) {
	utils.SetFormValues(form, topic, "Id")
	topic.CategoryId = int64(form.Category)
}

type CategoryAdminForm struct {
	Create bool   `form:"-"`
	Id     int    `form:"-"`
	Name   string `valid:"Required;MaxSize(30)"`
	Slug   string `valid:"Required;MaxSize(100)"`
	Order  int    ``
}

func (form *CategoryAdminForm) Labels() map[string]string {
	return map[string]string{
		"Name":  "model.category_name",
		"Slug":  "model.category_slug",
		"Order": "model.category_order",
	}
}

func (form *CategoryAdminForm) Valid(v *validation.Validation) {
	if models.IsExist(&models.Category{Name: form.Name, Id: int64(form.Id)}) {
		v.SetError("Name", "admin.field_need_unique")
	}

	if models.IsExist(&models.Category{Slug: form.Slug, Id: int64(form.Id)}) {
		v.SetError("Slug", "admin.field_need_unique")
	}
}

func (form *CategoryAdminForm) SetFromCategory(cat *models.Category) {
	utils.SetFormValues(cat, form)
}

func (form *CategoryAdminForm) SetToCategory(cat *models.Category) {
	utils.SetFormValues(form, cat, "Id")
}

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

package page

import (
	"github.com/astaxie/beego/validation"

	"github.com/go-tango/wego/modules/models"
	"github.com/go-tango/wego/modules/utils"
)

type PageAdminForm struct {
	Create     bool   `form:"-"`
	User       int    `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:"Required"`
	LastAuthor int    `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:""`
	Uri        string `valid:"Required;MaxSize(60);Match(/[0-9a-z-./]+/)"`
	Title      string `valid:"Required;MaxSize(60)"`
	Content    string `form:"type(textarea,markdown)" valid:"Required"`
	IsPublish  bool   ``
}

func (form *PageAdminForm) Valid(v *validation.Validation) {
	user := models.User{Id: form.User}
	if user.Read() != nil {
		v.SetError("User", "admin.not_found_by_id")
	}
}

func (form *PageAdminForm) SetFromPage(page *models.Page) {
	utils.SetFormValues(page, form)

	if page.User != nil {
		form.User = page.User.Id
	}

	if page.LastAuthor != nil {
		form.LastAuthor = page.LastAuthor.Id
	}
}

func (form *PageAdminForm) SetToPage(page *models.Page) {
	utils.SetFormValues(form, page)

	if page.User == nil {
		page.User = &models.User{}
	}
	page.User.Id = form.User

	if page.LastAuthor == nil {
		page.LastAuthor = &models.User{}
	}
	page.LastAuthor.Id = form.LastAuthor

	page.ContentCache = utils.RenderMarkdown(page.Content)
}

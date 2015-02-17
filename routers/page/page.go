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
	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/routers/base"
)

type Show struct {
	base.BaseRouter
}

func (this *Show) loadPage(page *models.Page) bool {
	uri := this.Ctx.Req().RequestURI

	page, err := models.GetPage(true, uri)
	if err == nil {
		this.Data["Page"] = page
	} else {
		this.NotFound()
	}
	return err != nil
}

func (this *Show) Get() {
	page := models.Page{}
	if this.loadPage(&page) {
		return
	}
	this.Render("page/show.html", this.Data)
}

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

	"github.com/go-tango/wetalk/modules/auth"
	"github.com/go-tango/wetalk/modules/utils"
	"github.com/go-tango/wetalk/routers/base"
)

type BaseAdminRouter struct {
	base.BaseRouter
}

func (this *BaseAdminRouter) NestPrepare() {
	if this.CheckActiveRedirect() {
		return
	}

	// if user isn't admin, then logout user
	if !this.User.IsAdmin {
		auth.LogoutUser(this.Context, this.Session.Session)
		// write flash message, use .flash.NotPermit
		this.FlashWrite("NotPermit", "true")
		this.Redirect("/login", 302)
		return
	}

	// it's admin and current in admin page
	this.Data["IsAdminPage"] = true

	/*if app, ok := this.AppController.(ModelPreparer); ok {
		app.ModelPrepare()
		return
	}*/
}

type ModelFinder interface {
	Object() interface{}
	ObjectQs() orm.QuerySeter
}

type ModelPreparer interface {
	ModelPrepare()
}

type ModelAdminRouter struct {
	BaseAdminRouter
}

func (this *ModelAdminRouter) ModelPrepare() {
	// set TplNames for model
	values := this.Ctx.Params()
	var tplNames string
	if model := values.Get(":model"); model != "" {
		if id := values.Get(":id"); id != "" {
			if id == "new" {
				tplNames = fmt.Sprintf("admin/%s/new.html", model)
			} else {
				if this.QueryObject() == false {
					return
				}

				if this.GetString(":action") == "delete" {
					tplNames = fmt.Sprintf("admin/%s/delete.html", model)
				} else {
					tplNames = fmt.Sprintf("admin/%s/edit.html", model)
				}
			}
		} else {
			tplNames = fmt.Sprintf("admin/%s/list.html", model)
		}

		name := fmt.Sprintf("%sAdmin", model)
		this.Data[name] = true
		this.RenderFile(tplNames, this.Data)
	}
}

// query objects and set to template
func (this *ModelAdminRouter) SetObjects(qs orm.QuerySeter, objects interface{}) error {
	cnt, err := qs.Count()
	if err != nil {
		return err
	}
	// create paginator
	p := this.SetPaginator(20, cnt)
	if cnt, err := qs.Limit(p.PerPageNums, p.Offset()).RelatedSel().All(objects); err != nil {
		return err
	} else {
		this.Data["Objects"] = objects
		this.Data["ObjectsCnt"] = cnt
	}
	return nil
}

// query object and set to template
func (this *ModelAdminRouter) QueryObject() bool {
	id, _ := utils.StrTo(this.GetString(":id")).Int()
	if id <= 0 {
		this.NotFound()
		return false
	}

	var app ModelFinder
	if a, ok := interface{}(this).(ModelFinder); ok {
		app = a
	} else {
		panic("ModelAdmin AppController need implement ModelFinder")
	}

	object := app.Object()
	qs := app.ObjectQs()

	// query object
	if err := qs.Filter("Id", id).Limit(1).One(object); err != nil {
		this.NotFound()
		if err != orm.ErrNoRows {
			beego.Error("SetObject: ", err)
		}
		return false

	} else {
		this.Data["Object"] = object
	}

	return true
}

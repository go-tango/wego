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
	"github.com/astaxie/beego/orm"
	"github.com/tango-contrib/xsrf"

	"github.com/go-tango/wego/models"
)

type ModelGet struct {
	BaseAdminRouter
	xsrf.NoCheck
}

func (this *ModelGet) Post() {
	this.Get()
}

func (this *ModelGet) Get() {
	id, _ := this.GetInt("id")
	model := this.GetString("model")
	result := map[string]interface{}{
		"success": false,
	}

	var data = make([][]interface{}, 0)

	defer func() {
		if len(data) > 0 {
			result["success"] = true
			result["data"] = data[0]
		}
		this.Data["json"] = result
		this.ServeJson(this.Data)
	}()

	if model == "User" {
		models.ORM().Iterate(&models.User{Id: id}, func(idx int, bean interface{}) error {
			user := bean.(*models.User)
			data = append(data, []interface{}{user.Id, user.UserName})
			return nil
		})
	}
}

type ModelSelect struct {
	BaseAdminRouter
	xsrf.NoCheck
}

func (this *ModelSelect) Post() {
	search := this.GetString("search")
	model := this.GetString("model")
	result := map[string]interface{}{
		"success": false,
	}

	var data []orm.ParamsList

	defer func() {
		if len(data) > 0 {
			result["success"] = true
			result["data"] = data
		}
		this.Data["json"] = result
		this.ServeJson(this.Data)
	}()

	if len(search) < 3 {
		return
	}

	if model == "User" {
		models.ORM().Limit(10).Where("user_name like ?", "%"+search+"%").
			Iterate(&models.User{}, func(idx int, bean interface{}) error {
			user := bean.(*models.User)
			data = append(data, []interface{}{user.Id, user.UserName})
			return nil
		})
	}
}

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

package api

import (
	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/auth"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/routers/base"
)

type Users struct {
	base.BaseRouter
}

func (this *Users) Post() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson(this.Data)
	}()

	if !this.IsAjax() {
		return
	}

	action := this.GetString("action")

	if this.IsLogin {

		switch action {
		case "get-follows":
			var data = make([][]interface{}, 0)
			models.ORM().Iterate(&models.Follow{UserId: this.User.Id},
				func(idx int, bean interface{}) error {
					followUser := bean.(*models.Follow).FollowUser()
					if followUser != nil {
						data = append(data, []interface{}{followUser.NickName, followUser.UserName})
					}
					return nil
				})
			result["success"] = true
			result["data"] = data

		case "follow", "unfollow":
			id, err := utils.StrTo(this.GetString("user")).Int()
			if err == nil && id != int(this.User.Id) {
				fuser := models.User{Id: int64(id)}
				if action == "follow" {
					auth.UserFollow(&this.User, &fuser)
				} else {
					auth.UserUnFollow(&this.User, &fuser)
				}
				result["success"] = true
			}
		}
	}
}

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

package auth

import (
	"github.com/astaxie/beego/orm"

	"github.com/go-tango/wego/modules/models"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/routers/base"
	"github.com/go-tango/wego/setting"
)

type UserRouter struct {
	base.BaseRouter
}

func (this *UserRouter) getUser(user *models.User) bool {
	username := this.Params().Get(":username")
	user.UserName = username

	err := user.Read("UserName")
	if err != nil {
		this.NotFound()
		return true
	}

	IsFollowed := false

	if this.IsLogin {
		if this.User.Id != user.Id {
			IsFollowed = this.User.FollowingUsers().Filter("FollowUser", user.Id).Exist()
		}
	}

	this.Data["TheUser"] = &user
	this.Data["IsFollowed"] = IsFollowed

	return false
}

type Home struct {
	UserRouter
}

func (this *Home) Get() {
	this.Data["IsUserHomePage"] = true

	var user models.User
	if this.getUser(&user) {
		return
	}

	//recent posts and comments
	limit := 5

	var posts []*models.Post
	var comments []*models.Comment

	user.RecentPosts().Limit(limit).RelatedSel().All(&posts)
	user.RecentComments().Limit(limit).RelatedSel().All(&comments)

	this.Data["TheUserPosts"] = posts
	this.Data["TheUserComments"] = comments

	//follow topics
	var ftopics []*models.FollowTopic
	var topics []*models.Topic
	ftNums, _ := models.FollowTopics().Filter("User", &user.Id).Limit(8).OrderBy("-Created").RelatedSel("Topic").All(&ftopics, "Topic")
	if ftNums > 0 {
		topics = make([]*models.Topic, 0, ftNums)
		for _, ft := range ftopics {
			topics = append(topics, ft.Topic)
		}
	}
	this.Data["TheUserFollowTopics"] = topics
	this.Data["TheUserFollowTopicsMore"] = ftNums >= 8

	//favorite posts
	var favPostIds orm.ParamsList
	var favPosts []models.Post
	favNums, _ := user.FavoritePosts().Limit(8).OrderBy("-Created").ValuesFlat(&favPostIds, "Post")
	if favNums > 0 {
		qs := models.Posts().Filter("Id__in", favPostIds)
		qs = qs.OrderBy("-Created").RelatedSel()
		models.ListObjects(qs, &favPosts)
	}
	this.Data["TheUserFavoritePosts"] = favPosts
	this.Data["TheUserFavoritePostsMore"] = favNums >= 8

	err := this.Render("user/home.html", this.Data)
	if err != nil {
		this.Result = err
	}
}

type Posts struct {
	UserRouter
}

func (this *Posts) Get() {
	var user models.User
	if this.getUser(&user) {
		return
	}

	limit := 20

	qs := user.RecentPosts()
	nums, _ := qs.Count()

	pager := this.SetPaginator(limit, nums)

	var posts []*models.Post
	qs.Limit(limit, pager.Offset()).RelatedSel().All(&posts)

	this.Data["TheUserPosts"] = posts
	this.Render("user/posts.html", this.Data)
}

type Comments struct {
	UserRouter
}

func (this *Comments) Get() {
	var user models.User
	if this.getUser(&user) {
		return
	}

	limit := 20

	qs := user.RecentComments()
	nums, _ := qs.Count()

	pager := this.SetPaginator(limit, nums)

	var comments []*models.Comment
	qs.Limit(limit, pager.Offset()).RelatedSel().All(&comments)

	this.Data["TheUserComments"] = comments

	this.Render("user/comments.html", this.Data)
}

func (this *UserRouter) getFollows(user *models.User, following bool) []map[string]interface{} {
	limit := 20

	var qs orm.QuerySeter

	if following {
		qs = user.FollowingUsers()
	} else {
		qs = user.FollowerUsers()
	}

	nums, _ := qs.Count()

	pager := this.SetPaginator(limit, nums)

	qs = qs.Limit(limit, pager.Offset())

	var follows []*models.Follow

	if following {
		qs.RelatedSel("FollowUser").All(&follows, "FollowUser")
	} else {
		qs.RelatedSel("User").All(&follows, "User")
	}

	if len(follows) == 0 {
		return nil
	}

	ids := make([]int, 0, len(follows))
	for _, follow := range follows {
		if following {
			ids = append(ids, follow.FollowUser.Id)
		} else {
			ids = append(ids, follow.User.Id)
		}
	}

	var eids orm.ParamsList
	this.User.FollowingUsers().Filter("FollowUser__in", ids).ValuesFlat(&eids, "FollowUser__Id")

	var fids map[int]bool
	if len(eids) > 0 {
		fids = make(map[int]bool)
		for _, id := range eids {
			tid, _ := utils.StrTo(utils.ToStr(id)).Int()
			if tid > 0 {
				fids[tid] = true
			}
		}
	}

	users := make([]map[string]interface{}, 0, len(follows))
	for _, follow := range follows {
		IsFollowed := false
		var u *models.User
		if following {
			u = follow.FollowUser
		} else {
			u = follow.User
		}
		if fids != nil {
			IsFollowed = fids[u.Id]
		}
		users = append(users, map[string]interface{}{
			"User":       u,
			"IsFollowed": IsFollowed,
		})
	}

	return users
}

type Following struct {
	UserRouter
}

func (this *Following) Get() {
	var user models.User
	if this.getUser(&user) {
		return
	}

	users := this.getFollows(&user, true)

	this.Data["TheUserFollowing"] = users
	this.Render("user/following.html", this.Data)
}

type Followers struct {
	UserRouter
}

func (this *Followers) Get() {
	var user models.User
	if this.getUser(&user) {
		return
	}

	users := this.getFollows(&user, false)

	this.Data["TheUserFollowers"] = users
	this.Render("user/followers.html", this.Data)
}

type FollowTopics struct {
	UserRouter
}

func (this *FollowTopics) Get() {
	this.TplNames = "user/follow-topics.html"

	var user models.User
	if this.getUser(&user) {
		return
	}

	var ftopics []*models.FollowTopic
	var topics []*models.Topic
	nums, _ := models.FollowTopics().Filter("User", &user.Id).OrderBy("-Created").RelatedSel("Topic").All(&ftopics, "Topic")
	if nums > 0 {
		topics = make([]*models.Topic, 0, nums)
		for _, ft := range ftopics {
			topics = append(topics, ft.Topic)
		}
	}
	this.Data["TheUserFollowTopics"] = topics
}

type FavoritePosts struct {
	UserRouter
}

func (this *FavoritePosts) Get() {
	this.TplNames = "user/favorite-posts.html"

	var user models.User
	if this.getUser(&user) {
		return
	}

	var postIds orm.ParamsList
	var posts []models.Post
	nums, _ := user.FavoritePosts().OrderBy("-Created").ValuesFlat(&postIds, "Post")
	if nums > 0 {
		qs := models.Posts().Filter("Id__in", postIds)
		cnt, _ := models.CountObjects(qs)
		pager := this.SetPaginator(setting.PostCountPerPage, cnt)
		qs = qs.OrderBy("-Created").Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()
		models.ListObjects(qs, &posts)
	}

	this.Data["TheUserFavoritePosts"] = posts
}

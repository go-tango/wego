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
	"github.com/go-tango/wego/models"
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

	u, err := models.GetUserByName(username)
	if err != nil {
		this.NotFound()
		return true
	}
	*user = *u

	IsFollowed := false

	if this.IsLogin {
		if this.User.Id != user.Id {
			IsFollowed = models.IsExist(&models.Follow{
				UserId:       this.User.Id,
				FollowUserId: user.Id,
			})
		}
	}

	this.Data["TheUser"] = &user
	this.Data["IsFollowed"] = IsFollowed

	return false
}

type Home struct {
	UserRouter
}

func (this *Home) Get() error {
	this.Data["IsUserHomePage"] = true

	var user models.User
	if this.getUser(&user) {
		return nil
	}

	//recent posts and comments
	limit := 5

	posts, _ := models.RecentPosts("recent", limit, 0)
	comments, _ := models.RecentCommentsByUserId(user.Id, limit)

	this.Data["TheUserPosts"] = posts
	this.Data["TheUserComments"] = comments

	//follow topics
	var topics []*models.Topic
	ftopics, _ := models.FindFollowTopic(user.Id, 8)
	if len(ftopics) > 0 {
		topics = make([]*models.Topic, 0, len(ftopics))
		for _, ft := range ftopics {
			topics = append(topics, ft.Topic())
		}
	}
	this.Data["TheUserFollowTopics"] = topics
	this.Data["TheUserFollowTopicsMore"] = len(ftopics) >= 8

	//favorite posts
	var favPostIds = make([]int64, 0)
	var favPosts []models.Post
	models.ORM().Limit(8).Desc("created").Iterate(new(models.FavoritePost), func(idx int, bean interface{}) error {
		favPostIds = append(favPostIds, bean.(*models.FavoritePost).PostId)
		return nil
	})
	if len(favPostIds) > 0 {
		models.ORM().In("id", favPostIds).Desc("created").Find(&favPosts)
	}
	this.Data["TheUserFavoritePosts"] = favPosts
	this.Data["TheUserFavoritePostsMore"] = len(favPostIds) >= 8

	return this.Render("user/home.html", this.Data)
}

type Posts struct {
	UserRouter
}

func (this *Posts) Get() error {
	var user models.User
	if this.getUser(&user) {
		return nil
	}

	limit := 20
	nums, _ := models.Count(&models.Post{UserId: user.Id})
	pager := this.SetPaginator(limit, nums)

	var posts = make([]*models.Post, 0)
	models.Find(limit, pager.Offset(), &posts)

	this.Data["TheUserPosts"] = posts
	return this.Render("user/posts.html", this.Data)
}

type Comments struct {
	UserRouter
}

func (this *Comments) Get() error {
	var user models.User
	if this.getUser(&user) {
		return nil
	}

	limit := 20
	nums, _ := models.CountCommentsByUserId(int64(user.Id))
	pager := this.SetPaginator(limit, nums)

	var comments = make([]*models.Comment, 0)
	models.Find(limit, pager.Offset(), &comments)

	this.Data["TheUserComments"] = comments

	return this.Render("user/comments.html", this.Data)
}

func (this *UserRouter) getFollows(user *models.User, following bool) []map[string]interface{} {
	var follow models.Follow
	if following {
		follow.UserId = user.Id
	} else {
		follow.FollowUserId = user.Id
	}

	nums, _ := models.Count(&follow)

	limit := 20
	pager := this.SetPaginator(limit, nums)

	var follows []*models.Follow
	models.ORM().Limit(limit, pager.Offset()).Find(&follows, &follow)

	if len(follows) == 0 {
		return nil
	}

	ids := make([]int, 0, len(follows))
	for _, follow := range follows {
		if following {
			ids = append(ids, int(follow.FollowUserId))
		} else {
			ids = append(ids, int(follow.UserId))
		}
	}

	var fids = make(map[int]bool)
	models.ORM().In("follow_user_id", ids).Iterate(&models.Follow{UserId: this.User.Id},
		func(idx int, bean interface{}) error {
			tid, _ := utils.StrTo(utils.ToStr(bean.(*models.Follow).Id)).Int()
			if tid > 0 {
				fids[tid] = true
			}
			return nil
		})

	users := make([]map[string]interface{}, 0, len(follows))
	for _, follow := range follows {
		IsFollowed := false
		var u *models.User
		if following {
			u = follow.FollowUser()
		} else {
			u = follow.User()
		}
		if fids != nil {
			IsFollowed = fids[int(u.Id)]
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

func (this *Following) Get() error {
	var user models.User
	if this.getUser(&user) {
		return nil
	}

	users := this.getFollows(&user, true)

	this.Data["TheUserFollowing"] = users
	return this.Render("user/following.html", this.Data)
}

type Followers struct {
	UserRouter
}

func (this *Followers) Get() error {
	var user models.User
	if this.getUser(&user) {
		return nil
	}

	users := this.getFollows(&user, false)

	this.Data["TheUserFollowers"] = users
	return this.Render("user/followers.html", this.Data)
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

	var topics []*models.Topic
	ftopics, _ := models.FindFollowTopic(user.Id, 0)
	if len(ftopics) > 0 {
		topics = make([]*models.Topic, 0, len(ftopics))
		for _, ft := range ftopics {
			topics = append(topics, ft.Topic())
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

	var postIds = make([]int64, 0)
	var posts []models.Post
	models.ORM().Desc("created").Iterate(&models.FavoritePost{UserId: user.Id},
		func(idx int, bean interface{}) error {
			postIds = append(postIds, bean.(*models.FavoritePost).PostId)
			return nil
		})
	if len(postIds) > 0 {
		cnt, _ := models.ORM().In("id", postIds).Count(models.Post{})
		pager := this.SetPaginator(setting.PostCountPerPage, cnt)
		models.ORM().Desc("created").Limit(setting.PostCountPerPage, pager.Offset()).Find(&posts)
	}

	this.Data["TheUserFavoritePosts"] = posts
}

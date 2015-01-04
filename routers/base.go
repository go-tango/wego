// Copyright 2013 wetalk authors
// Copyright 2014 wego authors
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

// An open source project for Gopher community.
package routers

import (
	"github.com/go-tango/wego/routers/admin"
	"github.com/go-tango/wego/routers/api"
	"github.com/go-tango/wego/routers/attachment"
	"github.com/go-tango/wego/routers/auth"
	"github.com/go-tango/wego/routers/base"
	"github.com/go-tango/wego/routers/page"
	"github.com/go-tango/wego/routers/post"
	"github.com/go-tango/wego/setting"

	"github.com/lunny/tango"
)

func Init(t *tango.Tango) {
	/* imgs */
	if setting.QiniuServiceEnabled {
		t.Get("/img/(.*)", attachment.QiniuImage)
	} else {
		t.Get("/img/(.*)", attachment.Image)
	}

	t.Use(setting.Captcha)
	// oauth support
	t.Get("/login/(.*)/access", new(auth.OAuthAccess))
	t.Get("/login/(.*)", new(auth.OAuthRedirect))

	/* Common Routers */
	t.Get("/", new(post.Home))
	t.Any("/topic/:slug", new(post.Topic))

	t.Get("/:sortSlug", new(post.Navs))
	t.Get("/category/:slug", new(post.Category))
	t.Get("/category/:catSlug/:sortSlug", new(post.CateNavs))

	t.Any("/new", new(post.NewPost))
	t.Any("/post/:post", new(post.SinglePost))
	t.Any("/post/:post/edit", new(post.EditPost))

	t.Get("/notification", new(post.NoticeRouter))

	if setting.SearchEnabled {
		t.Get("/search", new(post.SearchRouter))
	}

	t.Get("/user/:username/comments", new(auth.Comments))
	t.Get("/user/:username/posts", new(auth.Posts))
	t.Get("/user/:username/following", new(auth.Following))
	t.Get("/user/:username/followers", new(auth.Followers))
	t.Get("/user/:username/follow/topics", new(auth.FollowTopics))
	t.Get("/user/:username/favorite/posts", new(auth.FavoritePosts))
	t.Get("/user/:username", new(auth.Home))

	t.Any("/login", new(auth.Login))
	t.Get("/logout", new(auth.Logout))

	t.Any("/register/connect", new(auth.SocialAuthRouter))
	t.Any("/register", new(auth.Register))
	t.Get("/active/success", new(auth.RegisterSuccess))
	t.Get("/active/:code", new(auth.RegisterActive))

	t.Group("/settings", func(g *tango.Group) {
		g.Any("/profile", new(auth.ProfileRouter))
		g.Any("/change/password", new(auth.PasswordRouter))
		g.Any("/avatar", new(auth.AvatarRouter))
		g.Post("/avatar/upload", new(auth.AvatarUploadRouter))
	})

	t.Any("/forgot", new(auth.ForgotRouter))
	t.Any("/reset/:code", new(auth.ResetRouter))

	if setting.QiniuServiceEnabled {
		t.Post("/upload", new(attachment.QiniuUploadRouter))
	} else {
		t.Post("/upload", new(attachment.UploadRouter))
	}

	//download

	/* API Routers*/
	t.Post("/api/user", new(api.Users))
	t.Post("/api/md", new(api.Markdown))
	t.Post("/api/post", new(api.Post))

	// /* Admin Routers */
	t.Get("/admin", new(admin.AdminDashboard))

	t.Get("/admin/model/get", new(admin.ModelGet))
	t.Post("/admin/model/select", new(admin.ModelSelect))

	t.Get("/admin/user", new(admin.UserAdminList))
	t.Any("/admin/user/new", new(admin.UserAdminNew))
	t.Any("/admin/user/:id", new(admin.UserAdminEdit))
	t.Post("/admin/user/:id/:action", new(admin.UserAdminDelete))

	t.Get("/admin/post", new(admin.PostAdminList))
	t.Any("/admin/post/new", new(admin.PostAdminNew))
	t.Any("/admin/post/:id", new(admin.PostAdminEdit))
	t.Post("/admin/post/:id/:action", new(admin.PostAdminDelete))

	t.Get("/admin/comment", new(admin.CommentAdminList))
	t.Any("/admin/comment/new", new(admin.CommentAdminNew))
	t.Any("/admin/comment/:id", new(admin.CommentAdminEdit))
	t.Post("/admin/comment/:id/:action", new(admin.CommentAdminDelete))

	t.Get("/admin/topic", new(admin.TopicAdminList))
	t.Any("/admin/topic/new", new(admin.TopicAdminNew))
	t.Any("/admin/topic/:id", new(admin.TopicAdminEdit))
	t.Post("/admin/topic/:id/:action", new(admin.TopicAdminDelete))

	t.Get("/admin/category", new(admin.CategoryAdminList))
	t.Any("/admin/category/new", new(admin.CategoryAdminNew))
	t.Any("/admin/category/:id", new(admin.CategoryAdminEdit))
	t.Post("/admin/category/:id/:action", new(admin.CategoryAdminDelete))

	t.Get("/admin/page", new(admin.PageAdminList))
	t.Any("/admin/page/new", new(admin.PageAdminNew))
	t.Any("/admin/page/:id", new(admin.PageAdminEdit))
	t.Post("/admin/page/:id/:action", new(admin.PageAdminDelete))

	t.Get("/admin/bulletin", new(admin.BulletinAdminList))
	t.Any("/admin/bulletin/new", new(admin.BulletinAdminNew))
	t.Any("/admin/bulletin/:id", new(admin.BulletinAdminEdit))
	t.Post("/admin/bulletin/:id/:action", new(admin.BulletinAdminDelete))

	t.Get("/:slug", new(page.Show))

	// /* Robot routers for "robot.txt" */
	t.Get("/robot.txt", new(base.RobotRouter))
}

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
		t.Get("/img/*path", attachment.QiniuImage)
	} else {
		t.Get("/img/*path", attachment.Image)
	}

	// oauth support
	t.Get("/login/*auth/access", new(auth.OAuthAccess))
	t.Get("/login/*auth", new(auth.OAuthRedirect))

	/* Common Routers */
	t.Get("/", new(post.Home))
	t.Any("/topic/:slug", new(post.Topic))

	t.Get("/category/:slug", new(post.Category))
	t.Get("/category/:catSlug/:sortSlug", new(post.CateNavs))

	t.Any("/new", new(post.NewPost))
	t.Any("/post/:post", new(post.SinglePost))
	t.Any("/post/:post/edit", new(post.EditPost))

	t.Get("/notification", new(post.NoticeRouter))

	if setting.SearchEnabled {
		t.Get("/search", new(post.SearchRouter))
	}

	t.Group("/user/:username", func(g *tango.Group) {
		g.Get("/comments", new(auth.Comments))
		g.Get("/posts", new(auth.Posts))
		g.Get("/following", new(auth.Following))
		g.Get("/followers", new(auth.Followers))
		g.Get("/follow/topics", new(auth.FollowTopics))
		g.Get("/favorite/posts", new(auth.FavoritePosts))
		g.Get("", new(auth.Home))
	})

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

	/* API Routers*/
	t.Group("/api", func(g *tango.Group) {
		g.Post("/user", new(api.Users))
		g.Post("/md", new(api.Markdown))
		g.Post("/post", new(api.Post))
	})

	// /* Admin Routers */
	t.Group("/admin", func(g *tango.Group) {
		g.Get("", new(admin.AdminDashboard))
		g.Group("/model", func(cg *tango.Group) {
			cg.Any("/get", new(admin.ModelGet))
			cg.Post("/select", new(admin.ModelSelect))
		})

		g.Group("/user", func(cg *tango.Group) {
			cg.Get("", new(admin.UserAdminList))
			cg.Any("/new", new(admin.UserAdminNew))
			cg.Any("/:id", new(admin.UserAdminEdit))
			cg.Post("/:id/:action", new(admin.UserAdminDelete))
		})

		g.Group("/post", func(cg *tango.Group) {
			cg.Get("", new(admin.PostAdminList))
			cg.Any("/new", new(admin.PostAdminNew))
			cg.Any("/:id", new(admin.PostAdminEdit))
			cg.Post("/:id/:action", new(admin.PostAdminDelete))
		})

		g.Group("/comment", func(cg *tango.Group) {
			cg.Get("", new(admin.CommentAdminList))
			cg.Any("/new", new(admin.CommentAdminNew))
			cg.Any("/:id", new(admin.CommentAdminEdit))
			cg.Post("/:id/:action", new(admin.CommentAdminDelete))
		})

		g.Group("/topic", func(cg *tango.Group) {
			cg.Get("", new(admin.TopicAdminList))
			cg.Any("/new", new(admin.TopicAdminNew))
			cg.Any("/:id", new(admin.TopicAdminEdit))
			cg.Post("/:id/:action", new(admin.TopicAdminDelete))
		})

		g.Group("/category", func(cg *tango.Group) {
			cg.Get("", new(admin.CategoryAdminList))
			cg.Any("/new", new(admin.CategoryAdminNew))
			cg.Any("/:id", new(admin.CategoryAdminEdit))
			cg.Post("/:id/:action", new(admin.CategoryAdminDelete))
		})

		g.Group("/page", func(cg *tango.Group) {
			cg.Get("", new(admin.PageAdminList))
			cg.Any("/new", new(admin.PageAdminNew))
			cg.Any("/:id", new(admin.PageAdminEdit))
			cg.Post("/:id/:action", new(admin.PageAdminDelete))
		})

		g.Group("/bulletin", func(cg *tango.Group) {
			cg.Get("", new(admin.BulletinAdminList))
			cg.Any("/new", new(admin.BulletinAdminNew))
			cg.Any("/:id", new(admin.BulletinAdminEdit))
			cg.Post("/:id/:action", new(admin.BulletinAdminDelete))
		})
	})

	t.Get("/:sortSlug", new(post.Navs))
	t.Get("/page/:slug", new(page.Show))

	// /* Robot routers for "robot.txt" */
	t.Get("/robot.txt", new(base.RobotRouter))
}

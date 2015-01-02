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

// An open source project for Gopher community.
package routers

import (
	// "fmt"
	// "github.com/astaxie/beego"
	// "github.com/go-tango/wetalk/routers/admin"
	// "github.com/go-tango/wetalk/routers/api"
	// "github.com/go-tango/wetalk/routers/attachment"
	"github.com/go-tango/wetalk/routers/auth"
	// "github.com/go-tango/wetalk/routers/base"
	// "github.com/go-tango/wetalk/routers/page"
	 "github.com/go-tango/wetalk/routers/post"
	// "github.com/go-tango/wetalk/setting"

	"github.com/lunny/tango"
)

func Init(t *tango.Tango) {
	/* Add Filters */
	/*if setting.QiniuServiceEnabled {
		beego.InsertFilter("/img/*", beego.BeforeRouter, attachment.QiniuImageFilter)
	} else {
		beego.InsertFilter("/img/*", beego.BeforeRouter, attachment.ImageFilter)
	}

	beego.InsertFilter("/captcha/*", beego.BeforeRouter, setting.Captcha.Handler)
*/
	//beego.InsertFilter("/login/*/access", beego.BeforeRouter, auth.OAuthAccess)
	//beego.InsertFilter("/login/*", beego.BeforeRouter, auth.OAuthRedirect)

	/* Common Routers */
	t.Get("/", new(post.Home))
	t.Any("/topic/:slug", new(post.Topic))

	/*posts := new(post.PostListRouter)
	beego.Router("/:sortSlug(recent|hot|cold)", posts, "get:Navs")
	beego.Router("/category/:slug", posts, "get:Category")
	beego.Router("/category/:catSlug/:sortSlug(recent|hot|cold)", posts, "get:CatNavs")

	postR := new(post.PostRouter)
	beego.Router("/new", postR, "get:NewPost;post:NewPostSubmit")
	beego.Router("/post/:post([0-9]+)", postR, "get:SinglePost;post:SinglePostCommentSubmit")
	beego.Router("/post/:post([0-9]+)/edit", postR, "get:EditPost;post:EditPostSubmit")

	noticeRouter := new(post.NoticeRouter)
	beego.Router("/notification", noticeRouter, "get:Get")

	if setting.SearchEnabled {
		searchR := new(post.SearchRouter)
		beego.Router("/search", searchR, "get:Get")
	}

	user := new(auth.UserRouter)
	beego.Router("/user/:username/comments", user, "get:Comments")
	beego.Router("/user/:username/posts", user, "get:Posts")
	beego.Router("/user/:username/following", user, "get:Following")
	beego.Router("/user/:username/followers", user, "get:Followers")
	beego.Router("/user/:username/follow/topics", user, "get:FollowTopics")
	beego.Router("/user/:username/favorite/posts", user, "get:FavoritePosts")
	beego.Router("/user/:username", user, "get:Home")
*/

	t.Any("/login", new(auth.Login))
	t.Get("/logout", new(auth.Logout))

	//socialR := new(auth.SocialAuthRouter)
	//beego.Router("/register/connect", socialR, "get:Connect;post:ConnectPost")

	t.Any("/register", new(auth.Register))
	t.Get("/active/success", new(auth.RegisterSuccess))
	t.Get("/active/:code", new(auth.RegisterActive))

	/*
	settings := new(auth.SettingsRouter)
	beego.Router("/settings/profile", settings, "get:Profile;post:ProfileSave")
	beego.Router("/settings/change/password", settings, "get:ChangePassword;post:ChangePasswordSave")
	beego.Router("/settings/avatar", settings, "get:AvatarSetting;post:AvatarSettingSave")
	beego.Router("/settings/avatar/upload", settings, "post:AvatarUpload")

	forgot := new(auth.ForgotRouter)
	beego.Router("/forgot", forgot)
	beego.Router("/reset/:code([0-9a-zA-Z]+)", forgot, "get:Reset;post:ResetPost")

	if setting.QiniuServiceEnabled {
		upload := new(attachment.QiniuUploadRouter)
		beego.Router("/upload", upload, "post:Post")
	} else {
		upload := new(attachment.UploadRouter)
		beego.Router("/upload", upload, "post:Post")
	}

	//download

	/* API Routers*/
	// apiR := new(api.ApiRouter)
	// beego.Router("/api/user", apiR, "post:Users")
	// beego.Router("/api/md", apiR, "post:Markdown")
	// beego.Router("/api/post", apiR, "post:Post")

	// /* Admin Routers */
	// adminDashboard := new(admin.AdminDashboardRouter)
	// beego.Router("/admin", adminDashboard)

	// adminR := new(admin.AdminRouter)
	// beego.Router("/admin/model/get", adminR, "post:ModelGet")
	// beego.Router("/admin/model/select", adminR, "post:ModelSelect")

	// routes := map[string]beego.ControllerInterface{
	// 	"user":     new(admin.UserAdminRouter),
	// 	"post":     new(admin.PostAdminRouter),
	// 	"comment":  new(admin.CommentAdminRouter),
	// 	"topic":    new(admin.TopicAdminRouter),
	// 	"category": new(admin.CategoryAdminRouter),
	// 	"page":     new(admin.PageAdminRouter),
	// 	"bulletin": new(admin.BulletinAdminRouter),
	// }
	// for name, router := range routes {
	// 	beego.Router(fmt.Sprintf("/admin/:model(%s)", name), router, "get:List")
	// 	beego.Router(fmt.Sprintf("/admin/:model(%s)/:id(new)", name), router, "get:Create;post:Save")
	// 	beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)", name), router, "get:Edit;post:Update")
	// 	beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)/:action(delete)", name), router, "get:Confirm;post:Delete")
	// }
	// pageR := new(page.PageRouter)
	// beego.Router("/:slug", pageR, "get:Show")

	// /* Robot routers for "robot.txt" */
	// beego.Router("/robot.txt", &base.RobotRouter{})

}

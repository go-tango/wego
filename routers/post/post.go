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

package post

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/go-tango/wego/modules/models"
	"github.com/go-tango/wego/modules/post"
	"github.com/go-tango/wego/routers/base"
	"github.com/go-tango/wego/setting"

	"github.com/go-tango/wego/modules/utils"
	"github.com/lunny/log"
)

//Post List Router

type PostListRouter struct {
	base.BaseRouter
}

//Get all the categories
func (this *PostListRouter) setCategories(cats *[]models.Category) {
	//@see modules/post/topic_util.go
	post.ListCategories(cats)
	this.Data["Categories"] = *cats
}

//Get all the topics of the category
func (this *PostListRouter) setTopicsOfCategory(topics *[]models.Topic, category *models.Category) {
	//@see modules/post/topic_util.go
	post.ListTopicsOfCategory(topics, category)
	this.Data["TopicsOfCategory"] = *topics
}

//Get new best posts
func (this *PostListRouter) setNewBestPosts(posts *[]models.Post) {
	qs := models.Posts()
	qs = qs.Filter("IsBest", true).OrderBy("-Created").Limit(10)
	models.ListObjects(qs, posts)
	this.Data["NewBestPosts"] = posts
}

//Get new best posts by category
func (this *PostListRouter) setNewBestPostsOfCategory(posts *[]models.Post, cat *models.Category) {
	qs := models.Posts()
	qs = qs.Filter("IsBest", true).Filter("Category__id", cat.Id).OrderBy("-Created").Limit(10)
	models.ListObjects(qs, posts)
	this.Data["NewBestPosts"] = posts
}

//Get new best posts by topic
func (this *PostListRouter) setNewBestPostsOfTopic(posts *[]models.Post, topic *models.Topic) {
	qs := models.Posts()
	qs = qs.Filter("IsBest", true).Filter("Topic__id", topic.Id).OrderBy("-Created").Limit(10)
	models.ListObjects(qs, posts)
	this.Data["NewBestPosts"] = posts
}

//Get most replys posts
func (this *PostListRouter) setMostReplysPosts(posts *[]models.Post) {
	qs := models.Posts()
	qs = qs.Filter("Replys__gt", 0).OrderBy("-Created", "-Replys").Limit(10)
	models.ListObjects(qs, posts)
	this.Data["MostReplysPosts"] = posts
}

//Get most replys posts of category
func (this *PostListRouter) setMostReplysPostsOfCategory(posts *[]models.Post, cat *models.Category) {
	qs := models.Posts()
	qs = qs.Filter("Category__id", cat.Id).Filter("Replys__gt", 0).OrderBy("-Created", "-Replys").Limit(10)
	models.ListObjects(qs, posts)
	this.Data["MostReplysPosts"] = posts
}

//Get most replys post of topic
func (this *PostListRouter) setMostReplysPostsOfTopic(posts *[]models.Post, topic *models.Topic) {
	qs := models.Posts()
	qs = qs.Filter("Topic__id", topic.Id).Filter("Replys__gt", 0).OrderBy("-Created", "-Replys").Limit(10)
	models.ListObjects(qs, posts)
	this.Data["MostReplysPosts"] = posts
}

//Get sidebar bulletin information
func (this *PostListRouter) setSidebarBuilletinInfo() {
	var bulletins []models.Bulletin
	qs := models.Bulletins().OrderBy("Created")
	models.ListObjects(qs, &bulletins)

	var friendLinks []models.Bulletin
	var newComers []models.Bulletin
	var mobileApps []models.Bulletin
	var openSources []models.Bulletin

	for _, bulletin := range bulletins {
		switch bulletin.Type {
		case setting.BULLETIN_FRIEND_LINK:
			friendLinks = append(friendLinks, bulletin)
		case setting.BULLETIN_NEW_COMER:
			newComers = append(newComers, bulletin)
		case setting.BULLETIN_MOBILE_APP:
			mobileApps = append(mobileApps, bulletin)
		case setting.BULLETIN_OPEN_SOURCE:
			openSources = append(openSources, bulletin)
		}
	}
	this.Data["FriendLinks"] = friendLinks
	this.Data["NewComers"] = newComers
	this.Data["MobileApps"] = mobileApps
	this.Data["OpenSources"] = openSources
}

//Get the home page
type Home struct {
	PostListRouter
}

func (h *Home) Get() error {
	//get posts by Created datetime desc order
	var posts []models.Post
	qs := models.Posts()
	cnt, _ := models.CountObjects(qs)
	pager := h.SetPaginator(setting.PostCountPerPage, cnt)
	qs = qs.OrderBy("-LastReplied").Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()

	models.ListObjects(qs, &posts)
	h.Data["Posts"] = posts

	//top nav bar data
	var cats []models.Category
	h.setCategories(&cats)
	h.Data["SortSlug"] = ""
	h.Data["CategorySlug"] = "home"
	//new best posts
	var newBestPosts []models.Post
	h.setNewBestPosts(&newBestPosts)
	//most replys posts
	var mostReplysPosts []models.Post
	h.setMostReplysPosts(&mostReplysPosts)
	h.setSidebarBuilletinInfo()

	return h.Render("post/home.html", h.Data)
}

type Navs struct {
	PostListRouter
}

func (this *Navs) Get() {
	sortSlug := this.Params().Get(":sortSlug")
	var posts []models.Post
	qs := models.Posts()
	cnt, _ := models.CountObjects(qs)
	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	switch sortSlug {
	case "recent":
		qs = qs.OrderBy("-Created")
	case "hot":
		qs = qs.OrderBy("-LastReplied")
	case "cold":
		qs = qs.Filter("Replys", 0).OrderBy("-Created")
	default:
		this.NotFound()
		return
	}
	qs = qs.Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()
	models.ListObjects(qs, &posts)
	this.Data["Posts"] = posts

	//top nav bar data
	var cats []models.Category
	this.setCategories(&cats)
	this.Data["SortSlug"] = sortSlug
	this.Data["CategorySlug"] = "home"
	//new best posts
	var newBestPosts []models.Post
	this.setNewBestPosts(&newBestPosts)
	//most replys posts
	var mostReplysPosts []models.Post
	this.setMostReplysPosts(&mostReplysPosts)
	this.setSidebarBuilletinInfo()
	this.Render("post/home.html", this.Data)
}

type Category struct {
	PostListRouter
}

//Get the posts by category
func (this *Category) Get() {
	//check category slug
	slug := this.Params().Get(":slug")
	cat := models.Category{Slug: slug}
	if err := cat.Read("Slug"); err != nil {
		this.NotFound()
		return
	}
	//get posts by category slug, order by Created desc
	qs := models.Posts().Filter("Category", &cat)
	cnt, _ := models.CountObjects(qs)
	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	qs = qs.OrderBy("-LastReplied").Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()
	var posts []models.Post
	models.ListObjects(qs, &posts)

	this.Data["Category"] = &cat
	this.Data["Posts"] = posts

	//top nav bar data
	var cats []models.Category
	this.setCategories(&cats)
	var topics []models.Topic
	this.setTopicsOfCategory(&topics, &cat)
	this.Data["CategorySlug"] = cat.Slug
	this.Data["SortSlug"] = ""
	var newBestPosts []models.Post
	this.setNewBestPostsOfCategory(&newBestPosts, &cat)
	//most replys posts
	var mostReplysPosts []models.Post
	this.setMostReplysPostsOfCategory(&mostReplysPosts, &cat)
	this.setSidebarBuilletinInfo()
	this.Render("post/home.html", this.Data)
}

type CateNavs struct {
	PostListRouter
}

func (this *CateNavs) Get() {
	//check category slug and sort slug
	catSlug := this.Params().Get(":catSlug")
	sortSlug := this.Params().Get(":sortSlug")
	cat := models.Category{Slug: catSlug}
	if err := cat.Read("Slug"); err != nil {
		this.NotFound()
		return
	}
	qs := models.Posts().Filter("Category", &cat)
	cnt, _ := models.CountObjects(qs)
	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	switch sortSlug {
	case "recent":
		qs = qs.OrderBy("-Created")
	case "hot":
		qs = qs.OrderBy("-LastReplied")
	case "cold":
		qs = qs.Filter("Replys", 0).OrderBy("-Created")
	default:
		this.NotFound()
		return
	}
	qs = qs.Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()
	var posts []models.Post
	models.ListObjects(qs, &posts)

	this.Data["Category"] = &cat
	this.Data["Posts"] = posts

	//top nav bar data
	var cats []models.Category
	this.setCategories(&cats)
	var topics []models.Topic
	this.setTopicsOfCategory(&topics, &cat)
	this.Data["CategorySlug"] = cat.Slug
	this.Data["SortSlug"] = sortSlug
	var newBestPosts []models.Post
	this.setNewBestPostsOfCategory(&newBestPosts, &cat)
	//most replys posts
	var mostReplysPosts []models.Post
	this.setMostReplysPostsOfCategory(&mostReplysPosts, &cat)
	this.setSidebarBuilletinInfo()
	this.Render("post/home.html", this.Data)
}

type Topic struct {
	PostListRouter
}

//Topic Home Page
func (this *Topic) Get() {
	//check topic slug
	slug := this.Params().Get(":slug")
	topic := models.Topic{Slug: slug}
	if err := topic.Read("Slug"); err != nil {
		this.NotFound()
		return
	}
	//get topic category
	category := models.Category{Id: topic.Category.Id}
	if err := category.Read("Id"); err != nil {
		this.NotFound()
		return
	}

	//get posts by topic
	qs := models.Posts().Filter("Topic", &topic)
	cnt, _ := models.CountObjects(qs)
	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	qs = qs.OrderBy("-LastReplied").Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()
	var posts []models.Post
	models.ListObjects(qs, &posts)

	this.Data["Posts"] = posts
	this.Data["Topic"] = &topic
	this.Data["Category"] = &category

	//check whether added it into favorite list
	HasFavorite := false
	if this.IsLogin {
		HasFavorite = models.FollowTopics().Filter("User", &this.User).Filter("Topic", &topic).Exist()
	}
	this.Data["HasFavorite"] = HasFavorite

	//new best post
	var newBestPosts []models.Post
	this.setNewBestPostsOfTopic(&newBestPosts, &topic)
	//most replys posts
	var mostReplysPosts []models.Post
	this.setMostReplysPostsOfTopic(&mostReplysPosts, &topic)
	this.setSidebarBuilletinInfo()
	this.Render("post/topic.html", this.Data)
}

// Add this topic into favorite list
func (this *Topic) Post() {
	slug := this.Params().Get(":slug")

	topic := models.Topic{Slug: slug}
	if err := topic.Read("Slug"); err != nil {
		this.NotFound()
		return
	}

	result := map[string]interface{}{
		"success": false,
	}

	if this.IsAjax() {
		action := this.GetString("action")
		switch action {
		case "favorite":
			if this.IsLogin {
				qs := models.FollowTopics().Filter("User", &this.User).Filter("Topic", &topic)
				if qs.Exist() {
					qs.Delete()
				} else {
					fav := models.FollowTopic{User: &this.User, Topic: &topic}
					fav.Insert()
				}
				topic.RefreshFollowers()
				this.User.RefreshFavTopics()
				result["success"] = true
			}
		}
	}

	this.Data["json"] = result
	this.ServeJson(this.Data)
}

// Post Router
type PostRouter struct {
	base.BaseRouter
}

type NewPost struct {
	base.BaseRouter
}

func (this *NewPost) Get() {
	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}
	topicSlug := this.GetString("topic")
	if len(topicSlug) > 0 {
		topic := models.Topic{Slug: topicSlug}
		err := topic.Read("Slug")
		if err == nil {
			form.Topic = topic.Id
			form.Category = topic.Category.Id
			post.ListTopicsOfCategory(&form.Topics, &models.Category{Id: form.Category})
			this.Data["Topic"] = &topic
		} else {
			this.Redirect(setting.AppUrl)
			return
		}
	} else {
		catSlug := this.GetString("category")
		if len(catSlug) > 0 {
			category := models.Category{Slug: catSlug}
			category.Read("Slug")
			form.Category = category.Id
			post.ListTopicsOfCategory(&form.Topics, &category)
			this.Data["Category"] = &category
		} else {
			this.Redirect(setting.AppUrl)
			return
		}
	}

	this.SetFormSets(&form)
	this.Render("post/new.html", this.Data)
}

func (this *NewPost) Post() {
	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}
	topicSlug := this.GetString("topic")
	if len(topicSlug) > 0 {
		topic := models.Topic{Slug: topicSlug}
		err := topic.Read("Slug")
		if err == nil {
			form.Category = topic.Category.Id
			form.Topic = topic.Id
			this.Data["Topic"] = &topic
		} else {
			beego.Error("Can not find topic by slug:", topicSlug)
		}
	} else {
		topicId, err := this.GetInt("Topic")
		if err == nil {
			topic := models.Topic{Id: int(topicId)}
			err = topic.Read("Id")
			if err == nil {
				form.Category = topic.Category.Id
				form.Topic = topic.Id
				this.Data["Topic"] = &topic
			} else {
				log.Error("Can not find topic by id:", topicId)
			}
		} else {
			log.Error("Parse param Topic from request failed", err)
		}
	}
	if categorySlug := this.GetString("category"); categorySlug != "" {
		beego.Debug("Find category slug:", categorySlug)
		category := models.Category{Slug: categorySlug}
		category.Read("Slug")
		this.Data["Category"] = &category
	}
	post.ListTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	var post models.Post
	if err := form.SavePost(&post, &this.User); err == nil {
		this.JsStorage("deleteKey", "post/new")
		this.Redirect(post.Link())
		return
	}
	this.Render("post/new.html", this.Data)
}

func (this *PostRouter) loadPost(post *models.Post, user *models.User) bool {
	postId := this.Params().Get(":post")
	id, _ := strconv.ParseInt(postId, 10, 64)
	if id > 0 {
		qs := models.Posts().Filter("Id", id)
		if user != nil {
			qs = qs.Filter("User", user.Id)
		}
		qs.RelatedSel(1).One(post)
	}

	if post.Id == 0 {
		this.NotFound()
		return true
	}

	this.Data["Post"] = post

	return false
}

func (this *PostRouter) loadComments(post *models.Post, comments *[]*models.Comment) {
	qs := post.Comments()
	if num, err := qs.RelatedSel("User").OrderBy("Id").All(comments); err == nil {
		this.Data["Comments"] = *comments
		this.Data["CommentsNum"] = num
	}
}

type SinglePost struct {
	PostRouter
}

//Post Page
func (this *SinglePost) Get() {
	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return
	}

	var comments []*models.Comment
	this.loadComments(&postMd, &comments)

	//mark all notification as read
	if this.IsLogin {
		models.MarkNortificationAsRead(this.User.Id, postMd.Id)
	}

	//check whether this post is favorited
	num, _ := this.User.FavoritePosts().Filter("Post__Id", postMd.Id).Filter("IsFav", true).Count()
	if num != 0 {
		this.Data["IsPostFav"] = true
	} else {
		this.Data["IsPostFav"] = false
	}

	form := post.CommentForm{}
	this.SetFormSets(&form)
	//increment PageViewCount
	post.PostBrowsersAdd(this.User.Id, utils.IP(this.Req()), &postMd)
	this.Render("post/post.html", this.Data)
}

//New Comment
func (this *SinglePost) Post() {
	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return
	}

	var redir bool

	defer func() {
		if !redir {
			var comments []*models.Comment
			this.loadComments(&postMd, &comments)
		}
	}()

	form := post.CommentForm{}
	if !this.ValidFormSets(&form) {
		return
	}

	comment := models.Comment{}
	if err := form.SaveComment(&comment, &this.User, &postMd); err == nil {
		post.FilterCommentMentions(&this.User, &postMd, &comment)
		this.JsStorage("deleteKey", "post/comment")
		this.Redirect(postMd.Link(), 302)
		redir = true

		post.PostReplysCount(&postMd)
	}
	this.Render("post/post.html", this.Data)
}

type EditPost struct {
	PostRouter
}

func (this *EditPost) Get() {
	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	if !postMd.CanEdit {
		this.Redirect(postMd.Link())
		return
	}
	form := post.PostForm{}
	form.SetFromPost(&postMd)
	post.ListTopics(&form.Topics)
	this.SetFormSets(&form)
	this.Render("post/edit.html", this.Data)
}

func (this *PostRouter) Post() {
	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	if !postMd.CanEdit {
		this.FlashRedirect(postMd.Path(), 302, "CanNotEditPost")
	}

	form := post.PostForm{}
	form.SetFromPost(&postMd)
	post.ListTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	if err := form.UpdatePost(&postMd, &this.User); err == nil {
		this.JsStorage("deleteKey", "post/edit")
		this.Redirect(postMd.Link())
	}
	this.Render("post/edit.html", this.Data)
}

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

	"github.com/lunny/log"

	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/post"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/routers/base"
	"github.com/go-tango/wego/setting"
)

//Post List Router
type PostListRouter struct {
	base.BaseRouter
}

//Get all the categories
func (this *PostListRouter) setCategories(cats *[]models.Category) {
	models.FindCategories(cats)
	this.Data["Categories"] = *cats
}

func (this *PostListRouter) setTopics(topics *[]models.Topic) {
	models.FindTopics(topics)
	this.Data["Topics"] = *topics
}

//Get all the topics of the category
func (this *PostListRouter) setTopicsOfCategory(topics *[]models.Topic, category *models.Category) {
	models.FindTopicsByCategoryId(topics, category.Id)
	this.Data["TopicsOfCategory"] = *topics
}

//Get new best posts
func (this *PostListRouter) setNewBestPosts(posts *[]models.Post) {
	err := models.NewBestPostsByExample(posts, &models.Post{})
	if err != nil {
		this.Result = err
		return
	}

	this.Data["NewBestPosts"] = posts
}

//Get new best posts by category
func (this *PostListRouter) setNewBestPostsOfCategory(posts *[]models.Post, cat *models.Category) {
	err := models.NewBestPostsByExample(posts, &models.Post{CategoryId: cat.Id})
	if err != nil {
		this.Result = err
		return
	}

	this.Data["NewBestPosts"] = posts
}

//Get new best posts by topic
func (this *PostListRouter) setNewBestPostsOfTopic(posts *[]models.Post, topic *models.Topic) {
	err := models.NewBestPostsByExample(posts, &models.Post{TopicId: topic.Id})
	if err != nil {
		this.Result = err
		return
	}

	this.Data["NewBestPosts"] = posts
}

//Get most replys posts
func (this *PostListRouter) setMostReplysPosts(posts *[]models.Post) {
	err := models.MostReplysPostsByExample(posts, &models.Post{})
	if err != nil {
		this.Result = err
		return
	}
	this.Data["MostReplysPosts"] = posts
}

//Get most replys posts of category
func (this *PostListRouter) setMostReplysPostsOfCategory(posts *[]models.Post, cat *models.Category) {
	err := models.MostReplysPostsByExample(posts, &models.Post{CategoryId: cat.Id})
	if err != nil {
		this.Result = err
		return
	}
	this.Data["MostReplysPosts"] = posts
}

//Get most replys post of topic
func (this *PostListRouter) setMostReplysPostsOfTopic(posts *[]models.Post, topic *models.Topic) {
	err := models.MostReplysPostsByExample(posts, &models.Post{TopicId: topic.Id})
	if err != nil {
		this.Result = err
		return
	}
	this.Data["MostReplysPosts"] = posts
}

//Get sidebar bulletin information
func (this *PostListRouter) setSidebarBuilletinInfo() {
	bulletins, err := models.FindBulletins()
	if err != nil {
		this.Result = err
		return
	}

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
	cnt, err := models.CountByExample(&models.Post{})
	if err != nil {
		return err
	}

	pager := h.SetPaginator(setting.PostCountPerPage, cnt)
	posts, err := models.FindPosts(setting.PostCountPerPage, pager.Offset())
	if err != nil {
		return err
	}

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

func (this *Navs) Get() error {
	sortSlug := this.Params().Get(":sortSlug")

	cnt, err := models.CountByExample(&models.Post{})
	if err != nil {
		return err
	}

	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	posts, err := models.RecentPosts(sortSlug, setting.PostCountPerPage, pager.Offset())
	if err != nil {
		return err
	}

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
	return this.Render("post/home.html", this.Data)
}

type Category struct {
	PostListRouter
}

//Get the posts by category
func (this *Category) Get() error {
	//check category slug
	slug := this.Params().Get(":slug")
	cat, err := models.GetCategoryBySlug(slug)
	if err != nil {
		return err
	}

	//get posts by category slug, order by Created desc
	cnt, err := models.CountByExample(&models.Post{CategoryId: cat.Id})
	if err != nil {
		return err
	}

	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	posts, err := models.RecentPosts("hot", setting.PostCountPerPage, pager.Offset())
	if err != nil {
		return err
	}

	this.Data["Category"] = cat
	this.Data["Posts"] = posts

	//top nav bar data
	var cats []models.Category
	this.setCategories(&cats)
	var topics []models.Topic
	this.setTopicsOfCategory(&topics, cat)
	this.Data["CategorySlug"] = cat.Slug
	this.Data["SortSlug"] = ""
	var newBestPosts []models.Post
	this.setNewBestPostsOfCategory(&newBestPosts, cat)
	//most replys posts
	var mostReplysPosts []models.Post
	this.setMostReplysPostsOfCategory(&mostReplysPosts, cat)
	this.setSidebarBuilletinInfo()

	return this.Render("post/home.html", this.Data)
}

type CateNavs struct {
	PostListRouter
}

func (this *CateNavs) Get() error {
	//check category slug and sort slug
	catSlug := this.Params().Get(":catSlug")
	sortSlug := this.Params().Get(":sortSlug")
	cat, err := models.GetCategoryBySlug(catSlug)
	if err != nil {
		return err
	}

	cnt, err := models.CountByExample(&models.Post{CategoryId: cat.Id})
	if err != nil {
		return err
	}

	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	posts, err := models.RecentPosts(sortSlug, setting.PostCountPerPage, pager.Offset())
	if err != nil {
		return err
	}

	this.Data["Category"] = cat
	this.Data["Posts"] = posts

	//top nav bar data
	var cats []models.Category
	this.setCategories(&cats)
	var topics []models.Topic
	this.setTopicsOfCategory(&topics, cat)
	this.Data["CategorySlug"] = cat.Slug
	this.Data["SortSlug"] = sortSlug
	var newBestPosts []models.Post
	this.setNewBestPostsOfCategory(&newBestPosts, cat)
	//most replys posts
	var mostReplysPosts []models.Post
	this.setMostReplysPostsOfCategory(&mostReplysPosts, cat)
	this.setSidebarBuilletinInfo()
	return this.Render("post/home.html", this.Data)
}

type Topic struct {
	PostListRouter
}

//Topic Home Page
func (this *Topic) Get() error {
	//check topic slug
	slug := this.Params().Get(":slug")
	topic, err := models.GetTopicBySlug(slug)
	if err != nil {
		return err
	}

	//get topic category
	var category models.Category
	err = models.GetById(topic.CategoryId, &category)
	if err != nil {
		return err
	}

	//get posts by topic
	cnt, err := models.CountByExample(&models.Post{TopicId: topic.Id})
	if err != nil {
		return err
	}

	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	posts, err := models.FindPosts(setting.PostCountPerPage, pager.Offset())
	if err != nil {
		return err
	}

	this.Data["Posts"] = posts
	this.Data["Topic"] = &topic
	this.Data["Category"] = &category

	//check whether added it into favorite list
	var hasFavorite bool
	if this.IsLogin {
		hasFavorite, _ = models.HasUserFollowTopic(int64(this.User.Id), topic.Id)
	}
	this.Data["HasFavorite"] = hasFavorite

	//new best post
	var newBestPosts []models.Post
	this.setNewBestPostsOfTopic(&newBestPosts, topic)
	//most replys posts
	var mostReplysPosts []models.Post
	this.setMostReplysPostsOfTopic(&mostReplysPosts, topic)
	this.setSidebarBuilletinInfo()
	return this.Render("post/topic.html", this.Data)
}

// Add this topic into favorite list
func (this *Topic) Post() {
	slug := this.Params().Get(":slug")
	result := map[string]interface{}{
		"success": false,
	}

	topic, err := models.GetTopicBySlug(slug)
	if err != nil {
		return
	}

	if this.IsAjax() {
		action := this.GetString("action")
		switch action {
		case "favorite":
			if this.IsLogin {
				has, err := models.HasUserFollowTopic(int64(this.User.Id), topic.Id)
				if err != nil {
					log.Error("get follow user error:", err)
					return
				}

				if has {
					err = models.DeleteFollowTopic(int64(this.User.Id), topic.Id)
				} else {
					fav := models.FollowTopic{UserId: int64(this.User.Id), TopicId: topic.Id}
					err = models.Insert(&fav)
				}

				if err != nil {
					return
				}

				//TODO: add back
				//topic.RefreshFollowers()
				//this.User.RefreshFavTopics()
				result["success"] = true
			}
		}
	}

	this.Data["json"] = result
	this.ServeJson(this.Data)
}

type NewPost struct {
	base.BaseRouter
}

func (this *NewPost) Get() error {
	if this.CheckActiveRedirect() {
		return nil
	}

	form := post.PostForm{Locale: this.Locale}
	topicSlug := this.GetString("topic")
	if len(topicSlug) > 0 {
		topic, err := models.GetTopicBySlug(topicSlug)
		if err != nil {
			this.Redirect(setting.AppUrl)
			return nil
		}

		form.Topic = topic.Id
		form.Category = topic.CategoryId

		err = models.FindTopicsByCategoryId(&form.Topics, topic.CategoryId)
		if err != nil {
			this.Redirect(setting.AppUrl)
			return nil
		}

		this.Data["Topic"] = topic
	} else {
		catSlug := this.GetString("category")
		if len(catSlug) > 0 {
			category, err := models.GetCategoryBySlug(catSlug)
			if err != nil {
				return err
			}

			form.Category = category.Id
			err = models.FindTopicsByCategoryId(&form.Topics, category.Id)
			if err != nil {
				return err
			}
			this.Data["Category"] = category
		} else {
			this.Redirect(setting.AppUrl)
			return nil
		}
	}

	this.SetFormSets(&form)
	return this.Render("post/new.html", this.Data)
}

func (this *NewPost) Post() {
	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}
	topicSlug := this.GetString("topic")
	if len(topicSlug) > 0 {
		topic, err := models.GetTopicBySlug(topicSlug)
		if err == nil {
			form.Category = topic.CategoryId
			form.Topic = topic.Id
			this.Data["Topic"] = topic
		} else {
			log.Error("Can not find topic by slug:", topicSlug)
		}
	} else {
		topicId, err := this.GetInt("Topic")
		if err == nil {
			topic, err := models.GetTopicById(topicId)
			if err == nil {
				form.Category = topic.CategoryId
				form.Topic = topic.Id
				this.Data["Topic"] = topic
			} else {
				log.Error("Can not find topic by id:", topicId)
			}
		} else {
			log.Error("Parse param Topic from request failed", err)
		}
	}
	if categorySlug := this.GetString("category"); categorySlug != "" {
		log.Debug("Find category slug:", categorySlug)
		category, err := models.GetCategoryBySlug(categorySlug)
		if err != nil {
			log.Error("Get category error", err)
		}
		this.Data["Category"] = &category
	}
	models.FindTopics(&form.Topics)
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

// Post Router
type PostRouter struct {
	base.BaseRouter
}

func (this *PostRouter) loadPost(post *models.Post, user *models.User) bool {
	postId := this.Params().Get(":post")
	id, _ := strconv.ParseInt(postId, 10, 64)
	if id > 0 {
		var userId int64
		if user != nil {
			userId = user.Id
		}
		err := models.GetPost(id, userId, post)
		if err != nil {
			log.Error("loadPost error:", err)
			return true
		}
	}

	if post.Id == 0 {
		this.NotFound()
		return true
	}

	this.Data["Post"] = post

	return false
}

func (this *PostRouter) loadComments(post *models.Post, comments *[]*models.Comment) {
	err := models.GetCommentsByPostId(comments, post.Id)
	if err == nil {
		this.Data["Comments"] = *comments
		this.Data["CommentsNum"] = len(*comments)
	} else {
		log.Error("loadComments error:", err)
	}
}

type SinglePost struct {
	PostRouter
}

//Post Page
func (this *SinglePost) Get() error {
	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return nil
	}

	var comments []*models.Comment
	this.loadComments(&postMd, &comments)

	//mark all notification as read
	if this.IsLogin {
		models.MarkNortificationAsRead(this.User.Id, postMd.Id)
	}

	//check whether this post is favorited
	isPostFav, _ := models.IsPostFavorite(postMd.Id, int64(this.User.Id))
	this.Data["IsPostFav"] = isPostFav

	form := post.CommentForm{}
	this.SetFormSets(&form)
	//increment PageViewCount

	post.PostBrowsersAdd(this.User.Id, utils.IP(this.Req()), &postMd)
	return this.Render("post/post.html", this.Data)
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
	models.FindTopics(&form.Topics)
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
	models.FindTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	if err := form.UpdatePost(&postMd, &this.User); err == nil {
		this.JsStorage("deleteKey", "post/edit")
		this.Redirect(postMd.Link())
		return
	}
	this.Render("post/edit.html", this.Data)
}

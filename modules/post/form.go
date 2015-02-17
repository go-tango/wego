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
	"github.com/Unknwon/i18n"
	"github.com/astaxie/beego/validation"

	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
)

type PostForm struct {
	Lang     int            `form:"type(select);attr(rel,select2)"`
	Topic    int64          `form:"type(select);attr(rel,select2)" valid:"Required"`
	Title    string         `form:"attr(autocomplete,off)" valid:"Required;MinSize(5);MaxSize(60)"`
	Content  string         `form:"type(textarea)" valid:"Required;MinSize(10)"`
	Category int64          `form:"-"`
	Topics   []models.Topic `form:"-"`
	Locale   i18n.Locale    `form:"-"`
}

func (form *PostForm) LangSelectData() [][]string {
	langs := setting.Langs
	data := make([][]string, 0, len(langs))
	for i, lang := range langs {
		data = append(data, []string{lang, utils.ToStr(i)})
	}
	return data
}

func (form *PostForm) TopicSelectData() [][]string {
	data := make([][]string, 0, len(form.Topics))
	for _, topic := range form.Topics {
		data = append(data, []string{topic.Name, utils.ToStr(topic.Id)})
	}
	return data
}

func (form *PostForm) Valid(v *validation.Validation) {
	valid := false
	for _, topic := range form.Topics {
		if topic.Id == form.Topic {
			valid = true
		}
	}

	if !valid {
		v.SetError("Topic", "error")
	}

	if len(i18n.GetLangByIndex(form.Lang)) == 0 {
		v.SetError("Lang", "error")
	}
}

func (form *PostForm) SavePost(post *models.Post, user *models.User) error {
	utils.SetFormValues(form, post)
	post.CategoryId = form.Category
	post.TopicId = form.Topic
	post.UserId = user.Id
	post.LastReplyId = user.Id
	post.LastAuthorId = user.Id
	post.CanEdit = true
	post.ContentCache = utils.RenderMarkdown(form.Content)

	// mentioned follow users
	FilterMentions(user, post.ContentCache)

	return post.Insert()
}

func (form *PostForm) SetFromPost(post *models.Post) {
	utils.SetFormValues(post, form)
	form.Category = post.CategoryId
	form.Topic = post.TopicId
}

func (form *PostForm) UpdatePost(post *models.Post, user *models.User) error {
	changes := utils.FormChanges(post, form)
	if len(changes) == 0 {
		return nil
	}
	utils.SetFormValues(form, post)
	post.CategoryId = form.Category
	post.TopicId = form.Topic
	for _, c := range changes {
		if c == "Content" {
			post.ContentCache = utils.RenderMarkdown(form.Content)
			changes = append(changes, "ContentCache")
		}
	}

	// update last edit author
	if post.LastAuthorId != user.Id {
		post.LastAuthorId = user.Id
		changes = append(changes, "LastAuthor")
	}

	changes = append(changes, "Updated")

	return models.UpdateById(post.Id, post, models.Obj2Table(changes)...)
}

func (form *PostForm) Placeholders() map[string]string {
	return map[string]string{
		"Category": "model.category_choose_dot",
		"Topic":    "model.topic_choose_dot",
		"Title":    "post.plz_enter_title",
		"Content":  "post.plz_enter_content",
	}
}

type PostAdminForm struct {
	PostForm   `form:"-"`
	Create     bool   `form:"-"`
	User       int64  `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:"Required"`
	Title      string `valid:"Required;MaxSize(60)"`
	Content    string `form:"type(textarea,markdown)" valid:"Required"`
	Browsers   int    ``
	Replys     int    ``
	Favorites  int    ``
	LastReply  int64  `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:""`
	LastAuthor int64  `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:""`
	Topic      int64  `form:"type(select);attr(rel,select2)" valid:"Required"`
	Lang       int    `form:"type(select);attr(rel,select2)"`
	IsBest     bool   ``
}

func (form *PostAdminForm) Valid(v *validation.Validation) {
	var err error
	if _, err = models.GetUserById(form.User); err != nil {
		v.SetError("User", "admin.not_found_by_id")
	}

	if _, err = models.GetUserById(form.LastReply); err != nil {
		v.SetError("LastReply", "admin.not_found_by_id")
	}

	if _, err = models.GetUserById(form.LastAuthor); err != nil {
		v.SetError("LastReply", "admin.not_found_by_id")
	}

	if _, err = models.GetTopicById(form.Topic); err != nil {
		v.SetError("Topic", "admin.not_found_by_id")
	}

	if len(i18n.GetLangByIndex(form.Lang)) == 0 {
		v.SetError("Lang", "Not Found")
	}
}

func (form *PostAdminForm) SetFromPost(post *models.Post) {
	utils.SetFormValues(post, form)

	form.User = post.UserId
	form.LastReply = post.LastReplyId
	form.LastAuthor = post.LastAuthorId
	form.Topic = post.TopicId
}

func (form *PostAdminForm) SetToPost(post *models.Post) {
	utils.SetFormValues(form, post)

	post.UserId = form.User
	post.LastReplyId = form.LastReply
	post.LastAuthorId = form.LastAuthor
	post.TopicId = form.Topic
	//get category
	if topic, err := models.GetTopicById(form.Topic); err == nil {
		post.CategoryId = topic.CategoryId
	}
	post.ContentCache = utils.RenderMarkdown(post.Content)
}

type CommentForm struct {
	Message string `form:"type(textarea,markdown)" valid:"Required;MinSize(5)"`
}

func (form *CommentForm) SaveComment(comment *models.Comment, user *models.User, post *models.Post) error {
	comment.Message = form.Message
	comment.MessageCache = utils.RenderMarkdown(form.Message)
	comment.UserId = user.Id
	comment.PostId = post.Id
	if err := models.InsertComment(comment); err == nil {
		post.LastReplyId = user.Id
		models.UpdateById(post.Id, post, "last_reply_id", "last_replied")

		cnt, _ := models.CountCommentsLTEId(comment.Id)
		comment.Floor = int(cnt)
		return models.UpdateById(comment.Id, comment, "floor")
	} else {
		return err
	}
}

type CommentAdminForm struct {
	Create  bool   `form:"-"`
	User    int    `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:"Required"`
	Post    int    `valid:"Required"`
	Message string `form:"type(textarea)" valid:"Required"`
	Floor   int    `valid:"Required"`
	Status  int    `valid:""`
}

func (form *CommentAdminForm) Valid(v *validation.Validation) {
	var err error
	if _, err = models.GetUserById(int64(form.User)); err != nil {
		v.SetError("User", "admin.not_found_by_id")
	}

	if _, err = models.GetPostById(int64(form.Post)); err != nil {
		v.SetError("Post", "admin.not_found_by_id")
	}
}

func (form *CommentAdminForm) SetFromComment(comment *models.Comment) {
	utils.SetFormValues(comment, form)

	form.User = int(comment.UserId)
	form.Post = int(comment.PostId)
}

func (form *CommentAdminForm) SetToComment(comment *models.Comment) {
	utils.SetFormValues(form, comment)

	comment.UserId = int64(form.User)
	comment.PostId = int64(form.Post)
	comment.MessageCache = utils.RenderMarkdown(comment.Message)
}

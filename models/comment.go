package models

import (
	"time"

	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
)

// commnet content for post
type Comment struct {
	Id           int64
	UserId       int64  `xorm:"index"`
	PostId       int64  `xorm:"index"`
	Message      string `xorm:"text"`
	MessageCache string `xorm:"text"`
	Floor        int
	Status       int       `xorm:"index"`
	Created      time.Time `xorm:"created"`
}

func (m *Comment) GetMessageCache() string {
	if setting.RealtimeRenderMD {
		return utils.RenderMarkdown(m.Message)
	} else {
		return m.MessageCache
	}
}

func (m *Comment) String() string {
	return utils.ToStr(m.Id)
}

func (c *Comment) User() *User {
	var user User
	has, err := orm.Id(c.UserId).Get(&user)
	if err != nil || !has {
		return nil
	}
	return &user
}

func (c *Comment) Post() *Post {
	var post Post
	has, err := orm.Id(c.PostId).Get(&post)
	if err != nil || !has {
		return nil
	}
	return &post
}

func InsertComment(comment *Comment) error {
	_, err := orm.Insert(comment)
	return err
}

func RecentCommentsByUserId(userId int64, limit int) ([]Comment, error) {
	var comments = make([]Comment, 0)
	err := orm.Where("user_id = ?", userId).Limit(limit).Find(&comments)
	return comments, err
}

func GetCommentsByPostId(comments *[]*Comment, postId int64) error {
	return orm.Find(comments, &Comment{PostId: postId})
}

func CountCommentsByPostId(postId int64) (int64, error) {
	return orm.Count(&Comment{PostId: postId})
}

func CountCommentsByUserId(userId int64) (int64, error) {
	return orm.Count(&Comment{UserId: userId})
}

func CountCommentsLTEId(id int64) (int64, error) {
	return orm.Where("id <= ?", id).Count(new(Comment))
}

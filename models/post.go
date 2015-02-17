package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/Unknwon/i18n"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
)

// post content
type Post struct {
	Id           int64
	UserId       int64  `xorm:"index"`
	Title        string `xorm:"varchar(60)"`
	Content      string `xorm:"text"`
	ContentCache string `xorm:"text"`
	Browsers     int    `xorm:"index"`
	Replys       int    `xorm:"index"`
	Favorites    int    `xorm:"index"`
	LastReplyId  int64
	LastAuthorId int64
	TopicId      int64     `xorm:"index"`
	Lang         int       `xorm:"index"`
	IsBest       bool      `xorm:"index"`
	CanEdit      bool      `xorm:"index"`
	CategoryId   int64     `xorm:"index"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
	LastReplied  time.Time `xorm:"updated"`
}

func (m *Post) String() string {
	return utils.ToStr(m.Id)
}

func (m *Post) Link() string {
	return fmt.Sprintf("%spost/%d", setting.AppUrl, m.Id)
}

func (m *Post) Path() string {
	return fmt.Sprintf("/post/%d", m.Id)
}

func (m *Post) GetContentCache() string {
	if setting.RealtimeRenderMD {
		return utils.RenderMarkdown(m.Content)
	} else {
		return m.ContentCache
	}
}

func (m *Post) GetLang() string {
	return i18n.GetLangByIndex(m.Lang)
}

func (p *Post) Comments() []Comment {
	var comments = make([]Comment, 0)
	err := orm.Find(&comments, &Comment{PostId: p.Id})
	if err != nil {
		return nil
	}
	return comments
}

func (p *Post) Topic() *Topic {
	var topic Topic
	has, err := orm.Id(p.TopicId).Get(&topic)
	if err != nil || !has {
		return nil
	}
	return &topic
}

func (p *Post) Category() *Category {
	var category Category
	has, err := orm.Id(p.CategoryId).Get(&category)
	if err != nil || !has {
		return nil
	}
	return &category
}

func getUser(id int64) *User {
	var user User
	has, err := orm.Id(id).Get(&user)
	if err != nil || !has {
		return nil
	}
	return &user
}

func (p *Post) User() *User {
	return getUser(p.UserId)
}

func (p *Post) LastReply() *User {
	return getUser(p.LastReplyId)
}

func (p *Post) LastAuthor() *User {
	return getUser(p.LastAuthorId)
}

func (p *Post) Insert() error {
	_, err := orm.Insert(p)
	return err
}

func GetPost(id int64, userId int64, post *Post) error {
	s := orm.Where("id = ?", id)
	if userId > 0 {
		s.And("user_id = ?", userId)
	}
	has, err := s.Get(post)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

func GetPostById(id int64) (*Post, error) {
	var post Post
	has, err := orm.Id(id).Get(&post)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	return &post, nil
}

func FindPosts(limit, start int) ([]Post, error) {
	var posts = make([]Post, 0)
	err := orm.Desc("last_replied").Limit(limit, start).Find(&posts)
	return posts, err
}

func RecentPosts(sort string, limit, start int) ([]Post, error) {
	var posts = make([]Post, 0)
	s := orm.Limit(limit, start)
	switch sort {
	case "recent":
		s.Desc("created")
	case "hot":
		s.Desc("last_replied")
	case "cold":
		s.Where("Replys = ?", 0).Desc("created")
	default:
		return nil, errors.New("unknown sort")
	}
	err := s.Find(&posts)
	return posts, err
}

func NewBestPostsByExample(posts *[]Post, example *Post) error {
	return orm.Where("is_best = ?", true).Desc("created").Limit(10).Find(posts, example)
}

func MostReplysPostsByExample(posts *[]Post, example *Post) error {
	return orm.Where("replys > 0").Desc("created", "replys").Limit(10).Find(posts, example)
}

func UpdatePostBrowsersById(id int64) error {
	_, err := orm.Id(id).Incr("browsers").Update(new(Post))
	return err
}

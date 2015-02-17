package models

import (
	"fmt"
	"time"

	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
)

// post topic
type Topic struct {
	Id         int64
	Name       string    `xorm:"varchar(30) unique"`
	Intro      string    `xorm:"text"`
	ImageLink  string    `xorm:"varchar(200)"`
	Slug       string    `xorm:"varchar(100) unique"`
	Followers  int       `xorm:"index"`
	Order      int       `xorm:"index"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
	CategoryId int64     `xorm:"index"`
}

func (m *Topic) String() string {
	return utils.ToStr(m.Id)
}

func (m *Topic) Link() string {
	return fmt.Sprintf("%stopic/%s", setting.AppUrl, m.Slug)
}

func (t *Topic) Category() *Category {
	var category Category
	has, err := orm.Id(t.CategoryId).Get(&category)
	if err != nil || !has {
		return nil
	}
	return &category
}

func GetTopicByExample(topic *Topic) error {
	has, err := orm.Get(topic)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

func GetTopicById(id int64) (*Topic, error) {
	var topic = Topic{Id: id}
	err := GetTopicByExample(&topic)
	return &topic, err
}

func GetTopicBySlug(slug string) (*Topic, error) {
	var topic = Topic{Slug: slug}
	err := GetTopicByExample(&topic)
	return &topic, err
}

func FindTopics(topics *[]Topic) error {
	return orm.Find(topics)
}

func FindTopicsByCategoryId(topics *[]Topic, categoryId int64) error {
	return orm.Desc("order").Find(topics, &Topic{CategoryId: categoryId})
}

func CountTopicsByCategoryId(categoryId int64) (int64, error) {
	return orm.Count(&Topic{CategoryId: categoryId})
}

// user follow topics
type FollowTopic struct {
	Id      int64
	UserId  int64     `xorm:"unique(u)"`
	TopicId int64     `xorm:"unique(u)"`
	Created time.Time `xorm:"created"`
}

func (f *FollowTopic) Topic() *Topic {
	var topic Topic
	err := GetById(f.TopicId, &topic)
	if err != nil {
		return nil
	}
	return &topic
}

func FindFollowTopic(userId int64, limit int) ([]FollowTopic, error) {
	var follows = make([]FollowTopic, 0)
	sess := orm.Desc("created")
	if limit > 0 {
		sess.Limit(limit)
	}
	err := sess.Find(&follows, &FollowTopic{UserId: userId})
	return follows, err
}

func HasUserFollowTopic(userId, topicId int64) (bool, error) {
	has, err := orm.Get(&FollowTopic{UserId: userId, TopicId: topicId})
	if err != nil {
		return false, err
	}
	return has, nil
}

func DeleteFollowTopic(userId, topicId int64) error {
	_, err := orm.Delete(&FollowTopic{UserId: userId, TopicId: topicId})
	return err
}

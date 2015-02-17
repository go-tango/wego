package models

import "time"

type Page struct {
	Id           int64
	UserId       int64     `xorm:"index"`
	Uri          string    `xorm:"varchar(60) unqiue"`
	Title        string    `xorm:"varchar(60)"`
	Content      string    `xorm:"text"`
	ContentCache string    `xorm:"text"`
	LastAuthorId int64     `xorm:"index"`
	IsPublish    bool      `xorm:"index"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

func (p *Page) User() *User {
	return getUser(p.UserId)
}

func (p *Page) LastAuthor() *User {
	return getUser(p.LastAuthorId)
}

func GetPage(isPublish bool, uri string) (*Page, error) {
	var page = Page{
		IsPublish: isPublish,
		Uri:       uri,
	}
	has, err := orm.UseBool().Get(&page)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	return &page, nil
}

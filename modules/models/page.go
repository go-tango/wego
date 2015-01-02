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

package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/go-tango/wetalk/modules/utils"
	"github.com/go-tango/wetalk/setting"
)

type Page struct {
	Id           int
	User         *User     `orm:"rel(fk)"`
	Uri          string    `orm:"size(60);unqiue"`
	Title        string    `orm:"size(60)"`
	Content      string    `orm:"type(text)"`
	ContentCache string    `orm:"type(text)"`
	LastAuthor   *User     `orm:"rel(fk);null"`
	IsPublish    bool      `orm:"index"`
	Created      time.Time `orm:"auto_now_add"`
	Updated      time.Time `orm:"auto_now"`
}

func (m *Page) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Page) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Page) Update(fields ...string) error {
	fields = append(fields, "Updated")
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Page) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Page) String() string {
	return utils.ToStr(m.Id)
}

func (m *Page) Link() string {
	uri := m.Uri
	if len(uri) > 0 && uri[0] == '/' {
		uri = uri[1:]
	}
	return fmt.Sprintf("%s%s", setting.AppUrl, uri)
}

func (m *Page) GetTitle() string {
	return m.Title
}

func (m *Page) GetContentCache() string {
	var content, contentCache string
	content = m.Content
	contentCache = m.ContentCache

	if setting.RealtimeRenderMD {
		return utils.RenderMarkdown(content)
	} else {
		return contentCache
	}
}

func Pages() orm.QuerySeter {
	return orm.NewOrm().QueryTable("page").OrderBy("-Id")
}

func init() {
	orm.RegisterModel(new(Page))
}

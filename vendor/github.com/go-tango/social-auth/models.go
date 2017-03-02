// Copyright 2014 beego authors
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
//
// Maintain by https://github.com/slene

package social

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
)

const (
	startType SocialType = iota
	SocialGithub
	SocialGoogle
	SocialWeibo
	SocialQQ
	SocialDropbox
	SocialFacebook
	endType
)

var (
	types []SocialType
	orm   *xorm.Engine
)

func SetORM(o *xorm.Engine) {
	orm = o
}

func ORM() *xorm.Engine {
	return orm
}

func GetAllTypes() []SocialType {
	if types == nil {
		types = make([]SocialType, int(endType)-1)
		for i, _ := range types {
			types[i] = SocialType(i + 1)
		}
	}
	return types
}

type SocialType int

func (s SocialType) Available() bool {
	if s > startType && s < endType {
		return true
	}
	return false
}

func (s SocialType) Name() string {
	if p, ok := GetProviderByType(s); ok {
		return p.GetName()
	}
	return ""
}

func (s SocialType) NameLower() string {
	return strings.ToLower(s.Name())
}

type SocialTokenField struct {
	*Token
}

func (e *SocialTokenField) String() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *SocialTokenField) FromDB(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *SocialTokenField) ToDB() ([]byte, error) {
	return json.Marshal(e)
}

type UserSocial struct {
	Id       int64
	Uid      int              `xorm:"index"`
	Identify string           `xorm:"varchar(200)"`
	Type     SocialType       `xorm:"index"`
	Data     SocialTokenField `xorm:"text"`
}

func (e *UserSocial) Save() (err error) {
	if e.Id == 0 {
		_, err = orm.Insert(e)
	} else {
		_, err = orm.Id(e.Id).Update(e)
	}
	return
}

func (e *UserSocial) Token() (*Token, error) {
	return e.Data.Token, nil
}

func (e *UserSocial) PutToken(token *Token) error {
	if token == nil {
		return fmt.Errorf("token must be not nil")
	}

	changed := false

	if e.Data.Token == nil {
		e.Data.Token = token
		changed = true
	} else {

		if len(token.AccessToken) > 0 && token.AccessToken != e.Data.AccessToken {
			e.Data.AccessToken = token.AccessToken
			changed = true
		}
		if len(token.RefreshToken) > 0 && token.RefreshToken != e.Data.RefreshToken {
			e.Data.RefreshToken = token.RefreshToken
			changed = true
		}
		if len(token.TokenType) > 0 && token.TokenType != e.Data.TokenType {
			e.Data.TokenType = token.TokenType
			changed = true
		}
		if !token.Expiry.IsZero() && token.Expiry != e.Data.Expiry {
			e.Data.Expiry = token.Expiry
			changed = true
		}
	}

	if changed && e.Id > 0 {
		_, err := orm.Id(e.Id).Cols("data").Update(e)
		return err
	}

	return nil
}

func (e *UserSocial) TableUnique() [][]string {
	return [][]string{
		{"Identify", "Type"},
	}
}
func (e *UserSocial) Insert() error {
	_, err := orm.Insert(e)
	return err
}

/*
func (e *UserSocial) Read(fields ...string) error {
	if err := orm.NewOrm().Read(e, fields...); err != nil {
		return err
	}
	return nil
}*/

func Obj2Table(objs []string) []string {
	var res = make([]string, len(objs))
	for i, c := range objs {
		res[i] = orm.ColumnMapper.Obj2Table(c)
	}
	return res
}

func (e *UserSocial) Update(fields ...string) error {
	if _, err := orm.Id(e.Id).Cols(Obj2Table(fields)...).Update(e); err != nil {
		return err
	}
	return nil
}

func (e *UserSocial) Delete() error {
	if _, err := orm.Id(e.Id).Delete(new(UserSocial)); err != nil {
		return err
	}
	return nil
}

/*
func UserSocials() orm.QuerySeter {
	return orm.NewOrm().QueryTable("user_social")
}*/

// Get UserSocials by uid
func GetSocialsByUid(uid int, socialTypes ...SocialType) ([]*UserSocial, error) {
	var userSocials []*UserSocial
	err := orm.Where("uid = ?", uid).In("type", socialTypes).Find(&userSocials)
	if err != nil {
		return nil, err
	}
	return userSocials, nil
}

func init() {
	//orm.RegisterModel(new(UserSocial))
}

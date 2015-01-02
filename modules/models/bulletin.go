package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//Bulletin
type Bulletin struct {
	Id      int
	Name    string
	Url     string
	Type    int
	Created time.Time `orm:"auto_now_add"`
	Updated time.Time `orm:"auto_now"`
}

func (m *Bulletin) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Bulletin) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Bulletin) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Bulletin) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func Bulletins() orm.QuerySeter {
	return orm.NewOrm().QueryTable("Bulletin")
}

func init() {
	orm.RegisterModel(new(Bulletin))
}

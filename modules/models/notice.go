package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
	"time"
)

type Notification struct {
	Id           int
	FromUser     *User `orm:"rel(fk)"`
	ToUser       *User `orm:"rel(fk)"`
	Action       int
	Floor        int
	Lang         int
	TargetId     int
	Title        string    `orm:"size(60)"`
	Uri          string    `orm:"size(20)"`
	Content      string    `orm:"type(text)"`
	ContentCache string    `orm:"type(text)"`
	Status       int       `orm:"index"`
	Created      time.Time `orm:"auto_now_add"`
}

func GetUnreadNotificationCount(userId int) int64 {
	var qs = orm.NewOrm().QueryTable("notification")
	qs = qs.Filter("ToUser__Id", userId)
	qs = qs.Filter("Status", setting.NOTICE_UNREAD)
	qs = qs.Exclude("FromUser__Id", userId)
	if count, err := qs.Count(); err == nil {
		return count
	} else {
		return 0
	}
}

func MarkNortificationAsRead(userId int, postId int) {
	var query = fmt.Sprintf("UPDATE notification SET status=%d WHERE to_user_id=%d AND target_id=%d", setting.NOTICE_READ, userId, postId)
	var qs = orm.NewOrm().Raw(query)
	qs.Exec()
}

func Notifications(userId int) orm.QuerySeter {
	var qs = orm.NewOrm().QueryTable("notification")
	qs = qs.Filter("ToUser__Id", userId)
	qs = qs.Exclude("FromUser__Id", userId)
	qs = qs.OrderBy("-Id")
	return qs
}

func (m *Notification) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Notification) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Notification) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Notification) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Notification) String() string {
	return utils.ToStr(m.Id)
}

func (m *Notification) Link() string {
	return fmt.Sprintf("%s%s#reply%d", setting.AppUrl, m.Uri, m.Floor)
}

func (m *Notification) GetContentCache() string {
	if setting.RealtimeRenderMD {
		return utils.RenderMarkdown(m.Content)
	} else {
		return m.ContentCache
	}
}

func init() {
	orm.RegisterModel(new(Notification))
}

package models

import (
	"fmt"
	"time"

	"github.com/go-tango/wego/setting"
)

type Notification struct {
	Id           int64
	FromUserId   int64 `xorm:"index"`
	ToUserId     int64 `xorm:"index"`
	Action       int
	Floor        int
	Lang         int
	TargetId     int64
	Title        string    `xorm:"varchar(60)"`
	Uri          string    `xorm:"varchar(20)"`
	Content      string    `xorm:"text"`
	ContentCache string    `xorm:"text"`
	Status       int       `xorm:"index"`
	Created      time.Time `xorm:"created index"`
}

func (n *Notification) Link() string {
	return fmt.Sprintf("%notification", setting.AppUrl)
}

func (n *Notification) FromUser() *User {
	return getUser(n.FromUserId)
}

func (n *Notification) ToUser() *User {
	return getUser(n.ToUserId)
}

func InsertNotification(notic *Notification) error {
	_, err := orm.Insert(notic)
	return err
}

func CountNotifications(userId int64) (int64, error) {
	return orm.Where("from_user_id <> ?", userId).Count(&Notification{ToUserId: userId})
}

func FindNotificationsByUserId(userId int64, limit, start int) ([]*Notification, error) {
	var notifications = make([]*Notification, 0)
	err := orm.Desc("created").Limit(limit, start).Find(&notifications)
	return notifications, err
}

func MarkNortificationAsRead(userId int64, postId int64) error {
	_, err := orm.Exec("UPDATE notification SET status=? WHERE to_user_id=? AND target_id=?", setting.NOTICE_READ, userId, postId)
	return err
}

func GetUnreadNotificationCount(userId int64) int64 {
	count, _ := orm.Where("from_user_id <> ?", userId).Count(&Notification{
		ToUserId: userId,
		Status:   setting.NOTICE_UNREAD,
	})
	return count
}

package models

import (
	"fmt"
	"time"

	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
)

// main user table
// IsAdmin: user is admininstator
// IsActive: set active when email is verified
// IsForbid: forbid user login
type User struct {
	Id          int64
	UserName    string `xorm:"varchar(30) unique"`
	NickName    string `xorm:"varchar(30) unique"`
	Password    string `xorm:"varchar(128)"`
	AvatarType  int    `xorm:"default(1)"`
	AvatarKey   string `xorm:"varchar(50)"`
	Url         string `xorm:"varchar(100)"`
	Company     string `xorm:"varchar(30)"`
	Location    string `xorm:"varchar(30)"`
	Email       string `xorm:"varchar(80) unique"`
	GrEmail     string `xorm:"varchar(32)"`
	Info        string
	Github      string `xorm:"varchar(30)"`
	Twitter     string `xorm:"varchar(30)"`
	Google      string `xorm:"varchar(30)"`
	Weibo       string `xorm:"varchar(30)"`
	Linkedin    string `xorm:"varchar(30)"`
	Facebook    string `xorm:"varchar(30)"`
	PublicEmail bool
	Followers   int
	Following   int
	FavPosts    int
	FavTopics   int
	IsAdmin     bool      `xorm:"index"`
	IsActive    bool      `xorm:"index"`
	IsForbid    bool      `xorm:"index"`
	Lang        int       `xorm:"index"`
	Rands       string    `xorm:"varchar(10)"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
}

func (m *User) String() string {
	return utils.ToStr(m.Id)
}

func (m *User) Link() string {
	return fmt.Sprintf("%suser/%s", setting.AppUrl, m.UserName)
}

func (m *User) avatarLink(size int) string {
	if m.AvatarType == setting.AvatarTypePersonalized {
		if m.AvatarKey != "" {
			return fmt.Sprintf("%s", utils.GetQiniuZoomViewUrl(utils.GetQiniuPublicDownloadUrl(setting.QiniuAvatarDomain, m.AvatarKey), size, size))
		} else {
			return fmt.Sprintf("http://golanghome-public.qiniudn.com/golang_avatar.png?imageView/0/w/%s/h/%s/q/100", utils.ToStr(size), utils.ToStr(size))
		}
	} else {
		return fmt.Sprintf("%s%s?size=%s", setting.AvatarURL, m.GrEmail, utils.ToStr(size))
	}
}

func (m *User) AvatarLink24() string {
	return m.avatarLink(24)
}

func (m *User) AvatarLink48() string {
	return m.avatarLink(48)
}

func (m *User) AvatarLink64() string {
	return m.avatarLink(64)
}

func (m *User) AvatarLink100() string {
	return m.avatarLink(100)
}

func (m *User) AvatarLink200() string {
	return m.avatarLink(200)
}

func IsUserExistByName(username string, skipId int64) (bool, error) {
	return orm.Where("id <> ?", skipId).Get(&User{UserName: username})
}

func IsUserExistByEmail(email string, skipId int64) (bool, error) {
	return orm.Where("id <> ?", skipId).Get(&User{Email: email})
}

func GetUserById(id int64) (*User, error) {
	var user User
	has, err := orm.Id(id).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	return &user, nil
}

func GetUserByName(username string) (*User, error) {
	var user = User{
		UserName: username,
	}
	has, err := orm.Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	return &user, nil
}

// user follow
type Follow struct {
	Id           int64
	UserId       int64 `xorm:"index"`
	FollowUserId int64 `xorm:"index"`
	Mutual       bool
	Created      time.Time `orm:"created"`
}

func (f *Follow) FollowUser() *User {
	return getUser(f.FollowUserId)
}

func (f *Follow) User() *User {
	return getUser(f.UserId)
}

package models

import (
	"fmt"
	"time"

	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
)

type Image struct {
	Id      int64
	UserId  int64  `xorm:"index"`
	Token   string `xorm:"varchar(10)"`
	Width   int
	Height  int
	Ext     int       `xorm:"index"`
	Created time.Time `xorm:"created"`
}

func (m *Image) User() *User {
	return getUser(m.UserId)
}

func (m *Image) LinkFull() string {
	return m.LinkSize(0)
}

func (m *Image) LinkSmall() string {
	var width int
	switch {
	case m.Width > setting.ImageSizeSmall:
		width = setting.ImageSizeSmall
	}
	return m.LinkSize(width)
}

func (m *Image) LinkMiddle() string {
	var width int
	switch {
	case m.Width > setting.ImageSizeMiddle:
		width = setting.ImageSizeMiddle
	}
	return m.LinkSize(width)
}

func (m *Image) LinkSize(width int) string {
	if m.Ext == 3 {
		// if image is gif then return full size
		width = 0
	}
	var size string
	switch width {
	case setting.ImageSizeSmall, setting.ImageSizeMiddle:
		size = utils.ToStr(width)
	default:
		size = "full"
	}
	return "/img/" + m.GetToken() + "." + size + m.GetExt()
}

func (m *Image) GetExt() string {
	var ext string
	switch m.Ext {
	case 1:
		ext = ".jpg"
	case 2:
		ext = ".png"
	case 3:
		ext = ".gif"
	}
	return ext
}

func (m *Image) GetToken() string {
	number := utils.Date(m.Created, "ymds") + utils.ToStr(m.Id)
	return utils.NumberEncode(number, setting.ImageLinkAlphabets)
}

func (m *Image) DecodeToken(token string) error {
	number := utils.NumberDecode(token, setting.ImageLinkAlphabets)
	if len(number) < 9 {
		return fmt.Errorf("token `%s` too short <- `%s`", token, number)
	}

	if t, err := utils.DateParse(number[:8], "ymds"); err != nil {
		return fmt.Errorf("token `%s` date parse error <- `%s`", token, number)
	} else {
		m.Created = t
	}

	var err error
	m.Id, err = utils.StrTo(number[8:]).Int64()
	if err != nil {
		return fmt.Errorf("token `%s` id parse error <- `%s`", token, err)
	}

	return nil
}

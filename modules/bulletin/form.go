package bulletin

import (
	//"github.com/astaxie/beego/validation"
	"github.com/go-tango/wetalk/modules/models"
	"github.com/go-tango/wetalk/modules/utils"
	"github.com/go-tango/wetalk/setting"
)

type BulletinAdminForm struct {
	Create bool   `form:"-"`
	Id     int    `form:"-"`
	Name   string `valid:"Required;MaxSize(255);"`
	Url    string `valid:"Required;MaxSize(255);"`
	Type   int    `form:"type(select);attr(rel,select2)"`
}

func (form *BulletinAdminForm) TypeSelectData() [][]string {
	data := [][]string{
		[]string{"model.bulletin_friend_link", utils.ToStr(setting.BULLETIN_FRIEND_LINK)},
		[]string{"model.bulletin_new_comer", utils.ToStr(setting.BULLETIN_NEW_COMER)},
		[]string{"model.bulletin_mobile_app", utils.ToStr(setting.BULLETIN_MOBILE_APP)},
		[]string{"model.bulletin_open_source", utils.ToStr(setting.BULLETIN_OPEN_SOURCE)},
	}
	return data
}

func (form *BulletinAdminForm) Labels() map[string]string {
	return map[string]string{
		"Name": "model.bulletin_name",
		"Url":  "model.bulletin_url",
		"Type": "model.bulletin_type",
	}
}

func (form *BulletinAdminForm) SetFromBulletin(bulletin *models.Bulletin) {
	utils.SetFormValues(bulletin, form)
}

func (form *BulletinAdminForm) SetToBulletin(bulletin *models.Bulletin) {
	utils.SetFormValues(form, bulletin, "Id")
}

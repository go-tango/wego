package admin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/go-tango/wetalk/modules/bulletin"
	"github.com/go-tango/wetalk/modules/models"
	"github.com/go-tango/wetalk/modules/utils"
)

type BulletinAdminRouter struct {
	ModelAdminRouter
	object models.Bulletin
}

func (this *BulletinAdminRouter) Object() interface{} {
	return &this.object
}

func (this *BaseAdminRouter) ObjectQs() orm.QuerySeter {
	return models.Bulletins()
}
func (this *BulletinAdminRouter) List() {
	var bulletins []models.Bulletin
	qs := models.Bulletins().OrderBy("Type")
	if err := this.SetObjects(qs, &bulletins); err != nil {
		this.Data["Error"] = err
		beego.Error(err)
	}
}

func (this *BulletinAdminRouter) Create() {
	form := bulletin.BulletinAdminForm{Create: true}
	this.SetFormSets(&form)
}

func (this *BulletinAdminRouter) Save() {
	form := bulletin.BulletinAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var bulletin models.Bulletin
	form.SetToBulletin(&bulletin)
	if err := bulletin.Insert(); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/bulletin/%d", bulletin.Id), 302, "CreateSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}

func (this *BulletinAdminRouter) Edit() {
	form := bulletin.BulletinAdminForm{}
	form.SetFromBulletin(&this.object)
	this.SetFormSets(&form)
}
func (this *BulletinAdminRouter) Update() {
	form := bulletin.BulletinAdminForm{Id: this.object.Id}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/bulletin/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToBulletin(&this.object)
		if err := this.object.Update(changes...); err == nil {
			this.FlashRedirect(url, 302, "UpdateSuccess")
			return
		} else {
			beego.Error(err)
			this.Data["Error"] = err
		}
	} else {
		this.Redirect(url, 302)
	}
}
func (this *BulletinAdminRouter) Confirm() {

}
func (this *BulletinAdminRouter) Delete() {
	if this.FormOnceNotMatch() {
		return
	}
	qs := models.Bulletins().Filter("Id", this.object.Id)
	cnt, _ := qs.Count()
	if cnt > 0 {
		// delete object
		if err := this.object.Delete(); err == nil {
			this.FlashRedirect("/admin/bulletin", 302, "DeleteSuccess")
			return
		} else {
			beego.Error(err)
			this.Data["Error"] = err
		}
	}

}

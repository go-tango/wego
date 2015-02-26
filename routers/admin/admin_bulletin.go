package admin

import (
	"fmt"

	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/bulletin"
	"github.com/go-tango/wego/modules/utils"
	"github.com/lunny/log"
)

type BulletinAdminRouter struct {
	ModelAdminRouter
	object models.Bulletin
}

func (this *BulletinAdminRouter) Before() {
	this.Params().Set(":model", "bulletin")
	this.ModelAdminRouter.Before()
}

func (this *BulletinAdminRouter) Object() interface{} {
	return &this.object
}

type BulletinAdminList struct {
	BulletinAdminRouter
}

func (this *BulletinAdminList) Get() {
	var bulletins []models.Bulletin
	sess := models.ORM().Asc("type")
	if err := this.SetObjects(sess, &bulletins); err != nil {
		this.Data["Error"] = err
		log.Error(err)
	}
}

type BulletinAdminNew struct {
	BulletinAdminRouter
}

func (this *BulletinAdminNew) Get() {
	form := bulletin.BulletinAdminForm{Create: true}
	this.SetFormSets(&form)
}

func (this *BulletinAdminNew) Post() {
	form := bulletin.BulletinAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var bulletin models.Bulletin
	form.SetToBulletin(&bulletin)
	if err := models.Insert(&bulletin); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/bulletin/%d", bulletin.Id), 302, "CreateSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

type BulletinAdminEdit struct {
	BulletinAdminRouter
}

func (this *BulletinAdminEdit) Get() {
	form := bulletin.BulletinAdminForm{}
	form.SetFromBulletin(&this.object)
	this.SetFormSets(&form)
}

func (this *BulletinAdminEdit) Post() {
	form := bulletin.BulletinAdminForm{Id: int(this.object.Id)}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/bulletin/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToBulletin(&this.object)
		if err := models.UpdateById(this.object.Id, this.object, models.Obj2Table(changes)...); err == nil {
			this.FlashRedirect(url, 302, "UpdateSuccess")
			return
		} else {
			log.Error(err)
			this.Data["Error"] = err
		}
	} else {
		this.Redirect(url, 302)
	}
}

type BulletinAdminDelete struct {
	BulletinAdminRouter
}

func (this *BulletinAdminDelete) Post() {
	if this.FormOnceNotMatch() {
		return
	}
	cnt, _ := models.Count(&models.Bulletin{Id: this.object.Id})
	if cnt > 0 {
		// delete object
		if err := models.DeleteById(this.object.Id, new(models.Bulletin)); err == nil {
			this.FlashRedirect("/admin/bulletin", 302, "DeleteSuccess")
			return
		} else {
			log.Error(err)
			this.Data["Error"] = err
		}
	}

}

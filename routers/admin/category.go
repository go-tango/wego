package admin

import (
	"fmt"

	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/modules/post"
	"github.com/qiniu/log"
)

// view for new object save
func (this *CategoryAdminNew) Post() {
	form := post.CategoryAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var cat models.Category
	form.SetToCategory(&cat)
	if err := models.Insert(cat); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/category/%d", cat.Id), 302, "CreateSuccess")
		return
	} else {
		log.Error(err)
		this.Data["Error"] = err
	}
}

package post

import "github.com/go-tango/wego/models"

type NoticeRouter struct {
	PostListRouter
}

func (this *NoticeRouter) Get() error {
	this.Data["IsNotificationPage"] = true

	if this.CheckLoginRedirect() {
		return nil
	}

	pers := 10
	count, _ := models.CountNotifications(int64(this.User.Id))
	pager := this.SetPaginator(pers, count)

	notifications, err := models.FindNotificationsByUserId(int64(this.User.Id), pers, pager.Offset())
	if err != nil {
		return err
	}

	this.Data["Notifications"] = notifications

	var cats []models.Category
	var topics []models.Topic
	this.setCategories(&cats)
	this.setTopics(&topics)

	return this.Render("post/notice.html", this.Data)
}

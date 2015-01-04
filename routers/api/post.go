package api

import (
	"github.com/astaxie/beego/orm"
	"github.com/go-tango/wego/modules/models"
	"github.com/go-tango/wego/routers/base"
)

type Post struct {
	base.BaseRouter
}

func (this *Post) Post() {
	if this.CheckActiveRedirect() {
		return
	}

	if this.IsAjax() {
		result := map[string]interface{}{
			"success": false,
		}
		action := this.GetString("action")
		switch action {
		case "toggle-best":
			if !this.User.IsAdmin {
				result["success"] = false
			} else {
				if postId, err := this.GetInt("post"); err == nil {
					//set post best
					var post models.Post
					if err := orm.NewOrm().QueryTable("post").Filter("Id", postId).One(&post); err == nil {
						post.IsBest = !post.IsBest
						if post.Update("IsBest") == nil {
							result["success"] = true
						}
					}
				}
			}
		case "toggle-fav":
			if postId, err := this.GetInt("post"); err == nil {
				var post models.Post
				if err := orm.NewOrm().QueryTable("post").Filter("Id", postId).One(&post); err == nil {
					if post.Id != 0 {
						var favoritePost models.FavoritePost
						if this.User.FavoritePosts().Filter("Post__id", post.Id).One(&favoritePost); err == nil {
							if favoritePost.Id > 0 {
								//toogle IsFav
								favoritePost.IsFav = !favoritePost.IsFav
								if favoritePost.Update("IsFav") == nil {
									//update user fav post count
									if favoritePost.IsFav {
										this.User.FavPosts += 1
									} else {
										this.User.FavPosts -= 1

									}
									if this.User.Update("FavPosts") == nil {
										result["success"] = true
									}
								}
							} else {
								favoritePost = models.FavoritePost{
									User:  &this.User,
									Post:  &post,
									IsFav: true,
								}
								if favoritePost.Insert() == nil {
									//update user fav post count
									this.User.FavPosts += 1
									if this.User.Update("FavPosts") == nil {
										result["success"] = true
									}
								}
							}
						}
					}
				}
			}
		}
		this.Data["json"] = result
		this.ServeJson(this.Data)
	}
}

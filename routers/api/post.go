package api

import (
	"github.com/go-tango/wego/models"
	"github.com/go-tango/wego/routers/base"

	"github.com/tango-contrib/xsrf"
)

type Post struct {
	base.BaseRouter
	xsrf.NoCheck
}

func (this *Post) Post() {
	if this.CheckActiveRedirect() {
		return
	}

	if !this.IsAjax() {
		return
	}

	result := map[string]interface{}{
		"success": false,
	}
	action := this.GetString("action")
	switch action {
	case "toggle-best":
		if this.User.IsAdmin {
			if postId, err := this.GetInt("post"); err == nil {
				//set post best
				var post models.Post
				if err := models.GetById(postId, &post); err == nil {
					post.IsBest = !post.IsBest
					if models.UpdateById(post.Id, post, "is_best") == nil {
						result["success"] = true
					}
				}
			} else {
				this.Logger.Error("post value is not int:", this.GetString("post"))
			}
		}
	case "toggle-fav":
		if postId, err := this.GetInt("post"); err == nil {
			var post models.Post
			if err := models.GetById(postId, &post); err == nil {
				var favoritePost = models.FavoritePost{
					PostId: post.Id,
					UserId: this.User.Id,
				}

				if err := models.GetByExample(&favoritePost); err == nil {
					//toogle IsFav
					favoritePost.IsFav = !favoritePost.IsFav
					if models.UpdateById(favoritePost.Id, favoritePost, "is_fav") == nil {
						//update user fav post count
						if favoritePost.IsFav {
							this.User.FavPosts += 1
						} else {
							this.User.FavPosts -= 1

						}
						if models.UpdateById(this.User.Id, this.User, "fav_posts") == nil {
							result["success"] = true
						}
					}
				} else if err == models.ErrNotExist {
					favoritePost = models.FavoritePost{
						UserId: this.User.Id,
						PostId: post.Id,
						IsFav:  true,
					}
					if models.Insert(favoritePost) == nil {
						//update user fav post count
						this.User.FavPosts += 1
						if models.UpdateById(this.User.Id, this.User, "fav_posts") == nil {
							result["success"] = true
						}
					}
				} else {
					this.Logger.Error("Get favorite post err:", err)
				}
			}
		}
	}
	this.Data["json"] = result
	this.ServeJson(this.Data)
}

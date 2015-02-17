package models

import "time"

// user favorite posts
type FavoritePost struct {
	Id      int64
	UserId  int64 `xorm:"index"`
	PostId  int64 `xorm:"index"`
	IsFav   bool
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func (p *FavoritePost) User() *User {
	return getUser(p.UserId)
}

func (p *FavoritePost) Post() *Post {
	post, _ := GetPostById(p.PostId)
	return post
}

func IsPostFavorite(postId, userId int64) (bool, error) {
	var favorite = FavoritePost{
		UserId: userId,
		PostId: postId,
		IsFav:  true,
	}
	return orm.UseBool().Get(&favorite)
}

func FindFavoritesByUserId(userId int64) ([]*FavoritePost, error) {
	var favorites = make([]*FavoritePost, 0)
	err := orm.Desc("created").Find(&favorites, &FavoritePost{UserId: userId})
	return favorites, err
}

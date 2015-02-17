package models

import (
	"fmt"

	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"
)

// topic category
type Category struct {
	Id    int64
	Name  string `xorm:"varchar(30) unique"`
	Slug  string `xorm:"varchar(100) unique"`
	Order int    `xorm:"index"`
}

func (m *Category) String() string {
	return utils.ToStr(m.Id)
}

func (m *Category) Link() string {
	return fmt.Sprintf("%scategory/%s", setting.AppUrl, m.Slug)
}

func GetCategoryBySlug(slug string) (*Category, error) {
	var cate = Category{Slug: slug}
	err := GetByExample(&cate)
	return &cate, err
}

func CountCategoryBySlug(slug string) (int64, error) {
	return orm.Count(&Category{Slug: slug})
}

func FindCategories(cats *[]Category) (int64, error) {
	err := orm.Desc("order").Find(cats)
	return int64(len(*cats)), err
}

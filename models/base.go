package models

import "github.com/go-xorm/xorm"

func Orm() *xorm.Engine {
	return orm
}

func GetById(id int64, obj interface{}) error {
	has, err := orm.Id(id).Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

func GetByExample(obj interface{}) error {
	has, err := orm.Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

func CountByExample(obj interface{}) (int64, error) {
	return orm.Count(obj)
}

func Count(obj interface{}) (int64, error) {
	return orm.Count(obj)
}

func IsExist(obj interface{}) bool {
	has, _ := orm.Get(obj)
	return has
}

func Insert(obj interface{}) error {
	_, err := orm.Insert(obj)
	return err
}

func Find(limit, start int, objs interface{}) error {
	return orm.Limit(limit, start).Find(objs)
}

func DeleteById(id int64, obj interface{}) error {
	_, err := orm.Id(id).Delete(obj)
	return err
}

func Obj2Table(objs []string) []string {
	var res = make([]string, len(objs))
	for i, c := range objs {
		res[i] = orm.ColumnMapper.Obj2Table(c)
	}
	return res
}

func UpdateById(id int64, object interface{}, cols ...string) error {
	_, err := orm.Cols(cols...).Id(id).Update(object)
	return err
}

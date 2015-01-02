package models

import (
//	"github.com/go-xorm/xorm"
)

//var orm *xorm.Engine

func Init(isProMode bool) {
	// TODO: use xomr instead beego orm
	/*var err error
	orm, err = xorm.NewEngine(setting.DriverName, setting.DataSource)
	if err != nil {
		panic(err)
	}

	orm.SetMaxIdle(setting.MaxIdle)
	orm.SetMaxOpen(setting.MaxOpen)
	if !isProMode {
		orm.ShowSQL = true
	}*/
}

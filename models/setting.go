package models

// global settings name -> value
type Setting struct {
	Id      int
	Name    string `xorm:"varchar(100) unique"`
	Value   string `xorm:"text"`
	Updated string `xorm:"created"`
}

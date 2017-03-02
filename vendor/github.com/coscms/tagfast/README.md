tagfast
=======

golang：优化读取struct内的tag值（只解析一次，以后都从缓存中读取。官方的版本每次使用typ.Field(i).Tag.Get("tag1")都要解析一次，效率不高）

用法
=======
```
package main
import (
  "fmt"
  "reflect"
  "github.com/coscms/tagfast"
)

type Coscms struct {
  Url string `xorm:"not null default '' VARCHAR(255)" valid:"Requied" form_widget:"text"`
  Email string `xorm:"not null default '' VARCHAR(255)" valid:"Requied" form_widget:"text"`
}

func main(){
  m := Coscms{}
  t := reflect.TypeOf(m)
  for i := 0; i < t.NumField(); i++ {
    f := t.Field(i)
    widget:=tagfast.Tag(t,f,"form_widget")
    fmt.Println("widget:",widget)
    
    valid:=tagfast.Tag(t,f,"valid")
    fmt.Println("valid:",valid)
    
    xorm:=tagfast.Tag(t,f,"xorm")
    fmt.Println("xorm:",xorm)
    
  }
}
```

注意事项
=======
不要在同一个包内定义相同名称的结构体。
由于是按照“包名称+结构体名称”来缓存Tag，第二个结构体如果跟第一个结构体同名，
那么它们只会被当做同一个结构体来获取Tag，从而导致第二个具有相同名称的结构体获取Tag不正确。
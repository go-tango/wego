flash [![Build Status](https://drone.io/github.com/tango-contrib/flash/status.png)](https://drone.io/github.com/tango-contrib/flash/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/flash)](http://gocover.io/github.com/tango-contrib/flash)
======

Middleware flash is a tool for share data between requests for [Tango](https://github.com/lunny/tango). 

## Notice

This is a new version, it stores all data via [session](https://github.com/tango-contrib/session) not cookie. And it is slightly non-compitable with old version.

## Installation

    go get github.com/tango-contrib/flash

## Simple Example

```Go

import "github.com/tango-contrib/session"

type FlashAction struct {
    flash.Flash
}

func (x *FlashAction) Get() {
    x.Flash.Set("test", "test")
}

func (x *FlashAction) Post() {
   x.Flash.Get("test").(string) == "test"
}

func main() {
    t := tango.Classic()
    sessions := session.Sessions()
    t.Use(flash.Flashes(sessions))
    t.Any("/", new(FlashAction))
    t.Run()
}
```

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.
xsrf [![Build Status](https://drone.io/github.com/tango-contrib/xsrf/status.png)](https://drone.io/github.com/tango-contrib/xsrf/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/xsrf)](http://gocover.io/github.com/tango-contrib/xsrf)
======

Middleware xsrf is a xsrf checker for [Tango](https://github.com/lunny/tango). 

## Installation

    go get github.com/tango-contrib/xsrf

## Simple Example

```Go
type XsrfAction struct {
    render.Render
    xsrf.Checker
}

func (x *XsrfAction) Get() error {
    return x.Render("test.html", render.T{
        "XsrfFormHtml": x.XsrfFormHtml(),
    })
}

func (x *XsrfAction) Post() {
    // xsrf will be checked before this being called
}

func main() {
    t := tango.Classic()
    t.Use(xsrf.New(expireTime))
    t.Run()
}
```

If you don't want some action do not check, then
```Go
type NoCheckAction struct {
    xsrf.NoCheck
}

func (x *NoCheckAction) Post() {
    // xsrf will NOT be checked before this being called
}
```
will be ok.

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.
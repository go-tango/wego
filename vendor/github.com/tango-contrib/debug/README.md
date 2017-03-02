debug [![Build Status](https://drone.io/github.com/tango-contrib/debug/status.png)](https://drone.io/github.com/tango-contrib/debug/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/debug)](http://gocover.io/github.com/tango-contrib/debug)
======

Middleware debug is a debug middleware for [Tango](https://github.com/lunny/tango). 

## Installation

    go get github.com/tango-contrib/debug

## Simple Example

```Go
type DebugAction struct {
    tango.Ctx
}

func (c *DebugAction) Get() {
    c.Write([]byte("get"))
}

func main() {
    t := tango.Classic()
    t.Use(debug.Debug())
    t.Get("/", new(DebugAction))
    t.Run()
}
```

When you run this, then you will find debug info on console or log file, it will show you the request detail info and response detail.

```
[tango] 2015/03/04 06:44:06 [Debug] debug.go:53 [debug] request: GET http://localhost:3000/
[tango] 2015/03/04 06:44:06 [Debug] debug.go:55 [debug] head: map[]
[tango] 2015/03/04 06:44:06 [Debug] debug.go:66 [debug] ----------------------- end request
[tango] 2015/03/04 06:44:06 [Debug] debug.go:78 [debug] response ------------------ 200
[tango] 2015/03/04 06:44:06 [Debug] debug.go:80 [debug] head: map[]
[tango] 2015/03/04 06:44:06 [Debug] debug.go:83 [debug] body: debug
[tango] 2015/03/04 06:44:06 [Debug] debug.go:85 [debug] ----------------------- end response
```

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.
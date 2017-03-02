session [![Build Status](https://drone.io/github.com/tango-contrib/session/status.png)](https://drone.io/github.com/tango-contrib/session/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/session)](http://gocover.io/github.com/tango-contrib/session)
======

Session is a session middleware for [Tango](https://github.com/lunny/tango).

## Backend Supports

Currently session support some backends below:

* Memory - memory as a session store, this is the default store
* [nodb](http://github.com/tango-contrib/session-nodb) - nodb as a session store
* [redis](http://github.com/tango-contrib/session-redis) - redis server as a session store
* [ledis](http://github.com/tango-contrib/session-ledis) - ledis server as a session store
* [ssdb](http://github.com/tango-contrib/session-ssdb) - ssdb server as a session store

## Installation

    go get github.com/tango-contrib/session

## Simple Example

```Go
package main

import (
    "github.com/lunny/tango"
    "github.com/tango-contrib/session"
)

type SessionAction struct {
    session.Session
}

func (a *SessionAction) Get() string {
    a.Session.Set("test", "1")
    return a.Session.Get("test").(string)
}

func main() {
    o := tango.Classic()
    o.Use(session.New(session.Options{
        MaxAge:time.Minute * 20,
        }))
    o.Get("/", new(SessionAction))
}
```

## Getting Help

- [API Reference](https://gowalker.org/github.com/tango-contrib/session)

## License

This project is under BSD License. See the [LICENSE](LICENSE) file for the full license text.

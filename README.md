# WeGo

[简体中文](README_CN.md)

An open source forum system for Gopher community forked from WeTalk and rewritten via [tango](http://github.com/lunny/tango) & [xorm](http://xorm.io).

## Installation

### From source

```
go get github.com/go-tango/wego
cd $GOPATH/src/github.com/go-tango/wego
go build
```

Copy `conf/global/app.ini` to `conf/app.ini` and edit it. All configure has comment in it.

The files in `conf/` can overwrite `conf/global/` in runtime.

run `./wego` and then open `http://localhost:9000` in your web browser.

## Dependencies

Contrib

* Tango [https://github.com/lunny/tango](https://github.com/lunny/tango) (develop branch)
* Social-Auth [https://github.com/go-tango/social-auth](https://github.com/go-tango/social-auth)
* Compress [https://github.com/beego/compress](https://github.com/beego/compress)
* i18n [https://github.com/Unknwon/i18n](https://github.com/Unknwon/i18n)
* Mysql [https://github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
* goconfig [https://github.com/Unknwon/goconfig](https://github.com/Unknwon/goconfig)
* fsnotify [https://github.com/howeyc/fsnotify](https://github.com/howeyc/fsnotify)
* resize [https://github.com/nfnt/resize](https://github.com/nfnt/resize)
* blackfriday [https://github.com/slene/blackfriday](https://github.com/slene/blackfriday)

## Static Files

WeGo use `Google Closure Compile` and `Yui Compressor` compress js and css files.

So you could need Java Runtime. Or close this feature in code by yourself.

## Contact

Maintain by [lunny](https://github.com/lunny)

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

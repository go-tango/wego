# WeGo

An open source project for Gopher community forked from WeTalk.

### Usage

```
go get -u github.com/go-tango/wego
cd $GOPATH/src/github.com/go-tango/wego
```

I suggest you [update all Dependencies](#dependencies)

Copy `conf/global/app.ini` to `conf/app.ini` and edit it. All configure has comment in it.

The files in `conf/` can overwrite `conf/global/` in runtime.


**Run WeGo**

```
bee run watchall
```

### Dependencies

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

```
go get -u github.com/lunny/tango
cd $GOPATH/src/github.com/lunny/tango
```

Update all Dependencies

```
go get -u github.com/go-tango/social-auth
go get -u github.com/beego/compress
go get -u github.com/Unknwon/i18n
go get -u github.com/go-sql-driver/mysql
go get -u github.com/Unknwon/goconfig
go get -u github.com/howeyc/fsnotify
go get -u github.com/nfnt/resize
go get -u github.com/slene/blackfriday
```

### Static Files

WeGo use `Google Closure Compile` and `Yui Compressor` compress js and css files.

So you could need Java Runtime. Or close this feature in code by yourself.

### Contact

Maintain by [lunny](https://github.com/xiaolunwen)

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

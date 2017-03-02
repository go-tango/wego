# WeGo

WeGo是一个开源的论坛系统，最初是从WeTalk克隆而来，后使用 [tango](http://github.com/lunny/tango) 和 [xorm](http://xorm.io) 进行了重写。

## 安装

### 从源代码安装

```
go get github.com/go-tango/wego
cd $GOPATH/src/github.com/go-tango/wego
go build
```

如果需要自定义配置文件，请拷贝 `conf/global/app.ini` 到 `conf/app.ini` 然后进行修改。

在目录 `conf/` 里面的配置文件将会覆盖 `conf/global/` 里面的配置文件。

运行 `./wego`，然后在浏览器打开 `http://localhost:9000` 就可以看到论坛。

## 依赖

* Tango [https://github.com/lunny/tango](https://github.com/lunny/tango) (develop branch)
* Social-Auth [https://github.com/go-tango/social-auth](https://github.com/go-tango/social-auth)
* Compress [https://github.com/beego/compress](https://github.com/beego/compress)
* i18n [https://github.com/Unknwon/i18n](https://github.com/Unknwon/i18n)
* Mysql [https://github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
* goconfig [https://github.com/Unknwon/goconfig](https://github.com/Unknwon/goconfig)
* fsnotify [https://github.com/howeyc/fsnotify](https://github.com/howeyc/fsnotify)
* resize [https://github.com/nfnt/resize](https://github.com/nfnt/resize)
* blackfriday [https://github.com/slene/blackfriday](https://github.com/slene/blackfriday)

## 静态文件

WeGo 使用 `Google Closure Compile` 和 `Yui Compressor` 压缩 js 和 css 文件.

所以你可能需要Java运行时环境，或者您可以通过配置文件关闭此特性。

## 联系

此工程由 [lunny](https://github.com/lunny) 维护。

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

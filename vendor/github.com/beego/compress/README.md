# Beego Compress

Beego Compress provides an automated system for compressing JavaScript and Css files

It default use [Google Closure Compiler](https://code.google.com/p/closure-compiler/wiki/BinaryDownloads) for js, and [Yui Compressor](https://github.com/yui/yuicompressor/releases) for css

## Sample Usage with Beego

[After create a config file](#config-file), you can simple use it in beego.

Move **compiler.jar** and **yuicompressor.jar** to your beego app path. Parallel with static path.

BTW: Of course you can integrated it with other framework or use it as a command line tool.

```go
func SettingCompress() {
	// load json config file
	isProductMode := false
	setting, err := compress.LoadJsonConf("conf/compress.json", isProductMode, "http://127.0.0.1/")
	if err != nil {
		beego.Error(err)
		return
	}

	// after use this api, can run command from shell.
	setting.RunCommand()

	if isProductMode {
		// if in product mode, can use this api auto compress files
		setting.RunCompress(true, false, true)
	}

	// add func to FuncMap for template use
	beego.AddFuncMap("compress_js", setting.Js.CompressJs)
	beego.AddFuncMap("compress_css", setting.Css.CompressCss)
}
```

In tempalte usage

```html
...
<head>
	...
	{{compress_css "lib"}}
	{{compress_js "lib"}}
	{{compress_js "app"}}
</head>
...
```

#### Congratulations!! Let's see html results.

Render result when isProductMode is `false`

```html
<!-- Beego Compress group `lib` begin -->
<link rel="stylesheet" href="http://127.0.0.1/static_source/css/bootstrap.css?ver=1382331000" />
<link rel="stylesheet" href="http://127.0.0.1/static_source/css/bootstrap-theme.css?ver=1382322974" />
<link rel="stylesheet" href="http://127.0.0.1/static_source/css/font-awesome.min.css?ver=1378615042" />
<link rel="stylesheet" href="http://127.0.0.1/static_source/css/select2.css?ver=1382197742" />
<!-- end -->
<!-- Beego Compress group `lib` begin -->
<script type="text/javascript" src="http://127.0.0.1/static_source/js/jquery.min.js?ver=1378644427"></script>
<script type="text/javascript" src="http://127.0.0.1/static_source/js/bootstrap.js?ver=1382328826"></script>
<script type="text/javascript" src="http://127.0.0.1/static_source/js/lib.min.js?ver=1382328441"></script>
<script type="text/javascript" src="http://127.0.0.1/static_source/js/jStorage.js?ver=1382271840"></script>
<!-- end -->
<!-- Beego Compress group `app` begin -->
<script type="text/javascript" src="http://127.0.0.1/static_source/js/main.js?ver=1382195678"></script>
<script type="text/javascript" src="http://127.0.0.1/static_source/js/editor.js?ver=1382342779"></script>
<!-- end -->

```

Render result when isProductMode is `true`

```html
<link rel="stylesheet" href="http://127.0.0.1:8092/static/css/lib.min.css?ver=1382346563" />
<script type="text/javascript" src="http://127.0.0.1:8092/static/js/lib.min.js?ver=1382346557"></script>
<script type="text/javascript" src="http://127.0.0.1:8092/static/js/app.min.js?ver=1382346560"></script>
```

## Config file

Full config file example.

note: All json key are not case sensitive

**compress.json:**

```
{
	"Js": {
		// SrcPath is path of source file
		"SrcPath": "static_source/js",
		// DistPath is path of compressed file
		"DistPath": "static/js",
		// SrcURL is url prefix of source file
		"SrcURL": "static_source/js",
		// DistURL is url prefix of compressed file
		"DistURL": "static/js",
		"Groups": {
			// lib is the name of this compress group
			"lib": {
				// all compressed file will combined and save to DistFile
				"DistFile": "lib.min.js",
				// source files of this group
				"SourceFiles": [
					"jquery.min.js",
					"bootstrap.js",
					"lib.min.js",
					"jStorage.js"
				],
				// skip compress file list
				"SkipFiles": [
					"jquery.min.js",
					"lib.min.js"
				]
			},
			"app": {
				"DistFile": "app.min.js",
				"SourceFiles": [
					"main.js",
					"editor.js"
				]
			}
		}
	},
	"Css": {
		// config of css is same with js
		"SrcPath": "static_source/css",
		"DistPath": "static/css",
		"SrcURL": "static_source/css",
		"DistURL": "static/css",
		"Groups": {
			"lib": {
				"DistFile": "lib.min.css",
				"SourceFiles": [
					"bootstrap.css",
					"bootstrap-theme.css",
					"font-awesome.min.css",
					"select2.css"
				],
				"SkipFiles": [
					"font-awesome.min.css",
					"select2.css"
				]
			}
		}
	}
}
```

## Command mode

when use api `setting.RunCommand()`

```
$ go build app.go
$ ./app compress
compress command usage:

    js     - compress all js files
    css    - compress all css files
    all    - compress all files

$ ./app compress js -h
Usage of compress command: js:
  -force=false: force compress file
  -skip=false: skip all cached file
  -v=false: verbose info

$ ./app compress css -h
Usage of compress command: css:
  -force=false: force compress file
  -skip=false: skip all cached file
  -v=false: verbose info
```

```
use -force for re-create dist file but not skip cached.
use -skip can skip all cached file and re-compress.
```

## Custom Compress

All api can view in [GoWalker](http://gowalker.org/github.com/beego/compress)

* [TmpPath](http://gowalker.org/github.com/beego/compress#_variables) is default path of cached files.
* Can modify [JsFilters / CssFilters](http://gowalker.org/github.com/beego/compress#_variables) use your compress filter.
* Can modify [JsTagTemplate / CssTagTemplate]((http://gowalker.org/github.com/beego/compress#_variables)) with your `<script>` `<link>` tag.

##  Contact and Issue

All beego projects need your support.

Any suggestion are welcome, please [add new issue](https://github.com/beego/compress/issues/new) let me known.

## LICENSE

beego compress is licensed under the Apache Licence, Version 2.0 (http://www.apache.org/licenses/LICENSE-2.0.html).
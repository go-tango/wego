package compress

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

var (
	TmpPath           = "tmp"
	JsFilters         = []Filter{ClosureFilter}
	CssFilters        = []Filter{YuiFilter}
	JsTagTemplate, _  = template.New("").Parse(`<script type="text/javascript" src="{{.URL}}"></script>`)
	CssTagTemplate, _ = template.New("").Parse(`<link rel="stylesheet" href="{{.URL}}" />`)
)

type Filter func(source string) string

type JsCompresser interface {
	CompressJs(name string) template.HTML
	SetProMode(isPro bool)
	SetStaticURL(url string)
}

type CssCompresser interface {
	CompressCss(name string) template.HTML
	SetProMode(isPro bool)
	SetStaticURL(url string)
}

type Settings struct {
	Js  JsCompresser
	Css CssCompresser
}

func NewJsCompress(srcPath, distPath, srcURL, distURL string, groups map[string]Group) JsCompresser {
	compress := new(compressJs)
	compress.SrcPath = srcPath
	compress.DistPath = distPath
	compress.SrcURL = srcURL
	compress.DistURL = distURL
	compress.Groups = groups
	compress.StaticURL = "/"
	return compress
}

func NewCssCompress(srcPath, distPath, srcURL, distURL string, groups map[string]Group) CssCompresser {
	compress := new(compressCss)
	compress.SrcPath = srcPath
	compress.DistPath = distPath
	compress.SrcURL = srcURL
	compress.DistURL = distURL
	compress.Groups = groups
	compress.StaticURL = "/"
	return compress
}

func LoadJsonConf(filePath string, proMode bool, staticURL string) (setting *Settings, err error) {
	type Conf struct {
		Js  *compressJs
		Css *compressCss
	}

	var data []byte
	if file, err := os.Open(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Beego Compress: Conf Load %s", err.Error())
	} else {
		data, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("Beego Compress: Conf Read %s", err.Error())
		}
	}

	conf := Conf{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, fmt.Errorf("Beego Compress: Conf Parse %s", err.Error())
	}

	setting = new(Settings)
	if conf.Js != nil {
		setting.Js = conf.Js
	} else {
		setting.Js = new(compressJs)
	}

	if conf.Css != nil {
		setting.Css = conf.Css
	} else {
		setting.Css = new(compressCss)
	}

	if staticURL == "" {
		staticURL = "/"
	}

	setting.Js.SetProMode(proMode)
	setting.Css.SetProMode(proMode)

	setting.Js.SetStaticURL(staticURL)
	setting.Css.SetStaticURL(staticURL)

	return setting, nil
}

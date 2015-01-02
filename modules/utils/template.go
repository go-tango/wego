// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"time"
	"strings"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/Unknwon/i18n"

	"github.com/go-tango/wetalk/setting"
)

// get HTML i18n string
func i18nHTML(lang, format string, args ...interface{}) template.HTML {
	return template.HTML(i18n.Tr(lang, format, args...))
}

func boolicon(b bool) (s template.HTML) {
	if b {
		s = `<i style="color:green;" class="icon-check""></i>`
	} else {
		s = `<i class="icon-check-empty""></i>`
	}
	return
}

func date(t time.Time) string {
	return beego.Date(t, setting.DateFormat)
}

func datetime(t time.Time) string {
	return beego.Date(t, setting.DateTimeFormat)
}

func datetimes(t time.Time) string {
	return beego.Date(t, setting.DateTimeShortFormat)
}

func loadtimes(t time.Time) int {
	return int(time.Since(t).Nanoseconds() / 1e6)
}

func sum(base interface{}, value interface{}, params ...interface{}) (s string) {
	switch v := base.(type) {
	case string:
		s = v + ToStr(value)
		for _, p := range params {
			s += ToStr(p)
		}
	}
	return s
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func timesince(lang string, t time.Time) string {
	seconds := int(time.Since(t).Seconds())
	switch {
	case seconds < 60:
		return i18n.Tr(lang, "seconds_ago", seconds)
	case seconds < 60*60:
		return i18n.Tr(lang, "minutes_ago", seconds/60)
	case seconds < 60*60*24:
		return i18n.Tr(lang, "hours_ago", seconds/(60*60))
	case seconds < 60*60*24*100:
		return i18n.Tr(lang, "days_ago", seconds/(60*60*24))
	default:
		return beego.Date(t, setting.DateFormat)
	}
}

// create an login url with specify redirect to param
func loginto(uris ...string) template.HTMLAttr {
	var uri string
	if len(uris) > 0 {
		uri = uris[0]
	}
	to := fmt.Sprintf("%slogin", setting.AppUrl)
	if len(uri) > 0 {
		to += "?to=" + url.QueryEscape(uri)
	}
	return template.HTMLAttr(to)
}

// Substr returns the substr from start to length.
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

// Html2str returns escaping text convert from html.
func Html2str(html string) string {
	src := string(html)

	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//remove STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//remove SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

// DateFormat takes a time and a layout string and returns a string with the formatted date. Used by the template parser as "dateformat"
func DateFormat(t time.Time, layout string) (datestring string) {
	datestring = t.Format(layout)
	return
}

// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// Parse Date use PHP time format.
func DateParse(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

// Date takes a PHP like date func to Go's time format.
func Date(t time.Time, format string) string {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return t.Format(format)
}

// Compare is a quick and dirty comparison function. It will convert whatever you give it to strings and see if the two values are equal.
// Whitespace is trimmed. Used by the template parser as "eq".
func Compare(a, b interface{}) (equal bool) {
	equal = false
	if strings.TrimSpace(fmt.Sprintf("%v", a)) == strings.TrimSpace(fmt.Sprintf("%v", b)) {
		equal = true
	}
	return
}

/*func Config(returnType, key string, defaultVal interface{}) (value interface{}, err error) {
	switch returnType {
	case "String":
		value = AppConfig.String(key)
	case "Bool":
		value, err = AppConfig.Bool(key)
	case "Int":
		value, err = AppConfig.Int(key)
	case "Int64":
		value, err = AppConfig.Int64(key)
	case "Float":
		value, err = AppConfig.Float(key)
	case "DIY":
		value, err = AppConfig.DIY(key)
	default:
		err = errors.New("Config keys must be of type String, Bool, Int, Int64, Float, or DIY!")
	}

	if err != nil {
		if reflect.TypeOf(returnType) != reflect.TypeOf(defaultVal) {
			err = errors.New("defaultVal type does not match returnType!")
		} else {
			value, err = defaultVal, nil
		}
	} else if reflect.TypeOf(value).Kind() == reflect.String {
		if value == "" {
			if reflect.TypeOf(defaultVal).Kind() != reflect.String {
				err = errors.New("defaultVal type must be a String if the returnType is a String")
			} else {
				value = defaultVal.(string)
			}
		}
	}

	return
}*/

// Convert string to template.HTML type.
func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

// Htmlquote returns quoted html string.
func Htmlquote(src string) string {
	//HTML编码为实体符号
	/*
	   Encodes `text` for raw use in HTML.
	       >>> htmlquote("<'&\\">")
	       '&lt;&#39;&amp;&quot;&gt;'
	*/

	text := string(src)

	text = strings.Replace(text, "&", "&amp;", -1) // Must be done first!
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)
	text = strings.Replace(text, "'", "&#39;", -1)
	text = strings.Replace(text, "\"", "&quot;", -1)
	text = strings.Replace(text, "“", "&ldquo;", -1)
	text = strings.Replace(text, "”", "&rdquo;", -1)
	text = strings.Replace(text, " ", "&nbsp;", -1)

	return strings.TrimSpace(text)
}

// Htmlunquote returns unquoted html string.
func Htmlunquote(src string) string {
	//实体符号解释为HTML
	/*
	   Decodes `text` that's HTML quoted.
	       >>> htmlunquote('&lt;&#39;&amp;&quot;&gt;')
	       '<\\'&">'
	*/

	// strings.Replace(s, old, new, n)
	// 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换

	text := string(src)
	text = strings.Replace(text, "&nbsp;", " ", -1)
	text = strings.Replace(text, "&rdquo;", "”", -1)
	text = strings.Replace(text, "&ldquo;", "“", -1)
	text = strings.Replace(text, "&quot;", "\"", -1)
	text = strings.Replace(text, "&#39;", "'", -1)
	text = strings.Replace(text, "&gt;", ">", -1)
	text = strings.Replace(text, "&lt;", "<", -1)
	text = strings.Replace(text, "&amp;", "&", -1) // Must be done last!

	return strings.TrimSpace(text)
}

// UrlFor returns url string with another registered controller handler with params.
//	usage:
//
//	UrlFor(".index")
//	print UrlFor("index")
//  router /login
//	print UrlFor("login")
//	print UrlFor("login", "next","/"")
//  router /profile/:username
//	print UrlFor("profile", ":username","John Doe")
//	result:
//	/
//	/login
//	/login?next=/
//	/user/John%20Doe
//
//  more detail http://beego.me/docs/mvc/controller/urlbuilding.md
func UrlFor(endpoint string, values ...string) string {
	// TODO: fix this
	return ""
	//return BeeApp.Handlers.UrlFor(endpoint, values...)
}

// returns script tag with src string.
func AssetsJs(src string) template.HTML {
	text := string(src)

	text = "<script src=\"" + src + "\"></script>"

	return template.HTML(text)
}

// returns stylesheet link tag with src string.
func AssetsCss(src string) template.HTML {
	text := string(src)

	text = "<link href=\"" + src + "\" rel=\"stylesheet\" />"

	return template.HTML(text)
}

func FuncMap() template.FuncMap {
	// Register template functions.
	r := make(template.FuncMap)
	r["i18n"] = i18nHTML
	r["boolicon"] = boolicon
	r["date"] = date
	r["datetime"] = datetime
	r["datetimes"] = datetimes
	r["dict"] = dict
	r["timesince"] = timesince
	r["loadtimes"] = loadtimes
	r["sum"] = sum
	r["loginto"] = loginto
	r["isnotificationread"] = isnotificationread
	r["getbulletintype"] = getbulletintype
	r["dateformat"] = DateFormat
	r["date"] = Date
	r["compare"] = Compare
	r["substr"] = Substr
	r["html2str"] = Html2str
	r["str2html"] = Str2html
	r["htmlquote"] = Htmlquote
	r["htmlunquote"] = Htmlunquote
	r["urlfor"] = UrlFor
	//r["renderform"] = RenderForm
	r["assets_js"] = AssetsJs
	r["assets_css"] = AssetsCss
	return r
}

func RenderTemplate(TplNames string, Data map[string]interface{}) string {
	if beego.RunMode == "dev" {
		beego.BuildTemplate(beego.ViewsPath)
	}

	ibytes := bytes.NewBufferString("")
	if _, ok := beego.BeeTemplates[TplNames]; !ok {
		panic("can't find templatefile in the path:" + TplNames)
	}
	err := beego.BeeTemplates[TplNames].ExecuteTemplate(ibytes, TplNames, Data)
	if err != nil {
		beego.Trace("template Execute err:", err)
	}
	icontent, _ := ioutil.ReadAll(ibytes)
	return string(icontent)
}

func isnotificationread(status int) bool {
	var result = false
	if status == setting.NOTICE_READ {
		result = true
	}
	return result
}

func getbulletintype(lang string, t int) string {
	var typeStr string
	switch t {
	case setting.BULLETIN_FRIEND_LINK:
		typeStr = i18n.Tr(lang, "model.bulletin_friend_link")
	case setting.BULLETIN_MOBILE_APP:
		typeStr = i18n.Tr(lang, "model.bulletin_mobile_app")
	case setting.BULLETIN_NEW_COMER:
		typeStr = i18n.Tr(lang, "model.bulletin_new_comer")
	case setting.BULLETIN_OPEN_SOURCE:
		typeStr = i18n.Tr(lang, "model.bulletin_open_source")
	}
	return typeStr
}

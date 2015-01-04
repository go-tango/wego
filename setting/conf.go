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

// Package utils implemented some useful functions.

package setting

import (
	"fmt"
	"net/url"
	"os"
	"io"
	"path/filepath"
	"strings"
	"time"
	"html/template"

	"github.com/Unknwon/goconfig"
	"github.com/howeyc/fsnotify"

	"github.com/lunny/log"

	"github.com/macaron-contrib/cache"
	"github.com/astaxie/beego/orm"
	"github.com/tango-contrib/captcha"
	"github.com/beego/compress"
	"github.com/Unknwon/i18n"
	"github.com/go-tango/social-auth"
	"github.com/go-tango/social-auth/apps"
)

const (
	APP_VER = "0.1.0.0"
)

var (
	AppName             string
	AppVer              string
	AppHost             string
	AppUrl              string
	AppLogo             string
	EnforceRedirect     bool
	AvatarURL           string
	SecretKey           string
	IsProMode           bool
	ActiveCodeLives     int
	ResetPwdCodeLives   int
	DateFormat          string
	DateTimeFormat      string
	DateTimeShortFormat string
	TimeZone            string
	RealtimeRenderMD    bool
	ImageSizeSmall      int
	ImageSizeMiddle     int
	ImageLinkAlphabets  []byte
	ImageXSend          bool
	ImageXSendHeader    string
	Langs               []string

	LoginRememberDays int
	LoginMaxRetries   int
	LoginFailedBlocks int

	CookieRememberName string
	CookieUserName     string

	// search
	SearchEnabled bool

	// mail setting
	MailUser     string
	MailFrom     string
	MailHost     string
	MailAuthUser string
	MailAuthPass string
)

var (
	// OAuth
	GithubClientId     string
	GithubClientSecret string
	GoogleClientId     string
	GoogleClientSecret string
	WeiboClientId      string
	WeiboClientSecret  string
	QQClientId         string
	QQClientSecret     string
)

var (
	QiniuServiceEnabled bool
	QiniuAccessKey      string
	QiniuSecurityKey    string
	QiniuPostBucket     string
	QiniuPostDomain     string
	QiniuAvatarBucket   string
	QiniuAvatarDomain   string
)

var (
	PostCountPerPage int
)

const (
	LangEnUS = iota
	LangZhCN
)

const (
	BULLETIN_FRIEND_LINK = iota
	BULLETIN_NEW_COMER
	BULLETIN_OPEN_SOURCE
	BULLETIN_MOBILE_APP
)

const (
	AvatarImageMaxLength   = 500 * 1024
	AvatarTypeGravatar     = 1
	AvatarTypePersonalized = 2
)

const (
	NOTICE_TYPE_COMMENT   = 1
	NOTICE_TYPE_FAVOURITE = 2

	NOTICE_UNREAD = 1
	NOTICE_READ   = 2
)

var (
	// Social Auth
	GithubAuth *apps.Github
	GoogleAuth *apps.Google
	SocialAuth *social.SocialAuth
)

var (
	Cfg     *goconfig.ConfigFile
	Cache   cache.Cache
	Captcha *captcha.Captchas
)

var (
	GlobalConfPath   = "conf/global/app.ini"
	AppConfPath      = "conf/app.ini"
	CompressConfPath = "conf/compress.json"
)

var (
	DefaultLang = LangZhCN
)

var (
	EnableXSRF bool
	XSRFExpire int64
)

var (
	SessionProvider string
	SessionSavePath string
	SessionName string
	SessionCookieLifeTime int
	SessionGCMaxLifetime int64
)

var (
	DriverName string
	DataSource string
	MaxIdle int
	MaxOpen int
)

var (
	Log *log.Logger
)

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and log level.
	Cfg, err = goconfig.LoadConfigFile(GlobalConfPath)

	if Cfg == nil {
		Cfg, err = goconfig.LoadConfigFile(AppConfPath)
		if err != nil {
			fmt.Println("Fail to load configuration file: " + err.Error())
			os.Exit(2)
		}

	} else {
		Cfg.AppendFiles(AppConfPath)
	}

	Cfg.BlockMode = false

	// set time zone of wetalk system
	TimeZone = Cfg.MustValue("app", "time_zone", "UTC")
	if _, err := time.LoadLocation(TimeZone); err == nil {
		os.Setenv("TZ", TimeZone)
	} else {
		fmt.Println("Wrong time_zone: " + TimeZone + " " + err.Error())
		os.Exit(2)
	}

	// Trim 4th part.
	AppVer = strings.Join(strings.Split(APP_VER, ".")[:3], ".")

	IsProMode = Cfg.MustValue("app", "run_mode") == "pro"

	// cache system
	Cache, err = cache.NewCache("memory", `{"interval":360}`)
	Captcha = captcha.New(captcha.Options{}, Cache)
	Captcha.FieldIdName = "CaptchaId"
	Captcha.FieldCaptchaName = "Captcha"

	// session settings
	SessionProvider = Cfg.MustValue("session", "session_provider", "file")
	SessionSavePath = Cfg.MustValue("session", "session_path", "sessions")
	SessionName = Cfg.MustValue("session", "session_name", "wetalk_sess")
	SessionCookieLifeTime = Cfg.MustInt("session", "session_life_time", 0)
	SessionGCMaxLifetime = Cfg.MustInt64("session", "session_gc_time", 86400)

	EnableXSRF = true
	// xsrf token expire time
	XSRFExpire = 86400 * 365

	DriverName = Cfg.MustValue("orm", "driver_name", "mysql")
	DataSource = Cfg.MustValue("orm", "data_source", "root:@/wetalk?charset=utf8")
	MaxIdle = Cfg.MustInt("orm", "max_idle_conn", 30)
	MaxOpen = Cfg.MustInt("orm", "max_open_conn", 50)

	//set logger
	os.MkdirAll("./logs", os.ModePerm)
	f, err := os.Create("logs/wego.log")
	if err != nil {
		log.Panic("create log file failed:", err)
	}

	w := io.MultiWriter(f, os.Stdout)
	log.SetOutput(w)
	Log = log.Std

	if IsProMode {
		log.SetOutputLevel(log.Linfo)
	} else {
		log.SetOutputLevel(log.Ldebug)
	}

	// set default database
	err = orm.RegisterDataBase("default", DriverName, DataSource, MaxIdle, MaxOpen)
	if err != nil {
		log.Error(err)
	}
	orm.RunCommand()

	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		log.Error(err)
	}

	reloadConfig()

	social.DefaultAppUrl = AppUrl

	// OAuth
	var clientId, secret string

	clientId = Cfg.MustValue("oauth", "github_client_id", "your_client_id")
	secret = Cfg.MustValue("oauth", "github_client_secret", "your_client_secret")
	GithubAuth = apps.NewGithub(clientId, secret)

	clientId = Cfg.MustValue("oauth", "google_client_id", "your_client_id")
	secret = Cfg.MustValue("oauth", "google_client_secret", "your_client_secret")
	GoogleAuth = apps.NewGoogle(clientId, secret)

	err = social.RegisterProvider(GithubAuth)
	if err != nil {
		log.Error(err)
	}
	err = social.RegisterProvider(GoogleAuth)
	if err != nil {
		log.Error(err)
	}

	settingLocales()
	settingCompress()

	configWatcher()

	return Cfg
}

func reloadConfig() {
	AppName = Cfg.MustValue("app", "app_name", "WeTalk Community")

	AppHost = Cfg.MustValue("app", "app_host", "127.0.0.1:8092")
	AppUrl = Cfg.MustValue("app", "app_url", "http://127.0.0.1:8092/")
	AppLogo = Cfg.MustValue("app", "app_logo", "/static/img/logo.gif")
	AvatarURL = Cfg.MustValue("app", "avatar_url")

	EnforceRedirect = Cfg.MustBool("app", "enforce_redirect")

	DateFormat = Cfg.MustValue("app", "date_format")
	DateTimeFormat = Cfg.MustValue("app", "datetime_format")
	DateTimeShortFormat = Cfg.MustValue("app", "datetime_short_format")

	SecretKey = Cfg.MustValue("app", "secret_key")
	if len(SecretKey) == 0 {
		fmt.Println("Please set your secret_key in app.ini file")
	}

	ActiveCodeLives = Cfg.MustInt("app", "acitve_code_live_minutes", 180)
	ResetPwdCodeLives = Cfg.MustInt("app", "resetpwd_code_live_minutes", 180)

	LoginRememberDays = Cfg.MustInt("app", "login_remember_days", 7)
	LoginMaxRetries = Cfg.MustInt("app", "login_max_retries", 5)
	LoginFailedBlocks = Cfg.MustInt("app", "login_failed_blocks", 10)

	CookieRememberName = Cfg.MustValue("app", "cookie_remember_name", "wetalk_magic")
	CookieUserName = Cfg.MustValue("app", "cookie_user_name", "wetalk_powerful")

	RealtimeRenderMD = Cfg.MustBool("app", "realtime_render_markdown")

	ImageSizeSmall = Cfg.MustInt("image", "image_size_small")
	ImageSizeMiddle = Cfg.MustInt("image", "image_size_middle")

	if ImageSizeSmall <= 0 {
		ImageSizeSmall = 300
	}

	if ImageSizeMiddle <= ImageSizeSmall {
		ImageSizeMiddle = ImageSizeSmall + 400
	}

	str := Cfg.MustValue("image", "image_link_alphabets")
	if len(str) == 0 {
		str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	ImageLinkAlphabets = []byte(str)

	ImageXSend = Cfg.MustBool("image", "image_xsend", false)
	ImageXSendHeader = Cfg.MustValue("image", "image_xsend_header", "X-Accel-Redirect")

	MailUser = Cfg.MustValue("mailer", "mail_name", "WeTalk Community")
	MailFrom = Cfg.MustValue("mailer", "mail_from", "example@example.com")

	// set mailer connect args
	MailHost = Cfg.MustValue("mailer", "mail_host", "127.0.0.1:25")
	MailAuthUser = Cfg.MustValue("mailer", "mail_user", "example@example.com")
	MailAuthPass = Cfg.MustValue("mailer", "mail_pass", "******")

	// search setting
	SearchEnabled = Cfg.MustBool("search", "enabled")

	// OAuth
	GithubClientId = Cfg.MustValue("oauth", "github_client_id", "your_client_id")
	GithubClientSecret = Cfg.MustValue("oauth", "github_client_secret", "your_client_secret")
	GoogleClientId = Cfg.MustValue("oauth", "google_client_id", "your_client_id")
	GoogleClientSecret = Cfg.MustValue("oauth", "google_client_secret", "your_client_secret")
	WeiboClientId = Cfg.MustValue("oauth", "weibo_client_id", "your_client_id")
	WeiboClientSecret = Cfg.MustValue("oauth", "weibo_client_secret", "your_client_secret")
	QQClientId = Cfg.MustValue("oauth", "qq_client_id", "your_client_id")
	QQClientSecret = Cfg.MustValue("oauth", "qq_client_secret", "your_client_secret")

	//Qiniu
	QiniuServiceEnabled = Cfg.MustBool("qiniu", "qiniu_service_enabled", false)
	QiniuAccessKey = Cfg.MustValue("qiniu", "qiniu_access_key")
	QiniuSecurityKey = Cfg.MustValue("qiniu", "qiniu_security_key")
	QiniuPostBucket = Cfg.MustValue("qiniu", "qiniu_post_bucket")
	QiniuPostDomain = Cfg.MustValue("qiniu", "qiniu_post_domain")
	QiniuAvatarBucket = Cfg.MustValue("qiniu", "qiniu_avatar_bucket")
	QiniuAvatarDomain = Cfg.MustValue("qiniu", "qiniu_avatar_domain")

	//post
	PostCountPerPage = Cfg.MustInt("post", "post_count_per_page", 20)
}

func settingLocales() {
	// load locales with locale_LANG.ini files
	langs := "en-US|zh-CN"
	for _, lang := range strings.Split(langs, "|") {
		lang = strings.TrimSpace(lang)
		files := []string{"conf/" + "locale_" + lang + ".ini"}
		if fh, err := os.Open(files[0]); err == nil {
			fh.Close()
		} else {
			files = nil
		}
		if err := i18n.SetMessage(lang, "conf/global/"+"locale_"+lang+".ini", files...); err != nil {
			log.Error("Fail to set message file: " + err.Error())
			os.Exit(2)
		}
	}
	Langs = i18n.ListLangs()
}

var Funcs = make(template.FuncMap)

func settingCompress() {
	setting, err := compress.LoadJsonConf(CompressConfPath, IsProMode, AppUrl)
	if err != nil {
		log.Error(err)
		return
	}

	setting.RunCommand()

	if IsProMode {
		setting.RunCompress(true, false, true)
	}

	Funcs["compress_js"] = setting.Js.CompressJs
	Funcs["compress_css"] = setting.Css.CompressCss
}

var eventTime = make(map[string]int64)

func configWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("Failed start app watcher: " + err.Error())
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				switch filepath.Ext(event.Name) {
				case ".ini":
					if checkEventTime(event.Name) {
						continue
					}
					log.Info(event)

					if err := Cfg.Reload(); err != nil {
						log.Error("Conf Reload: ", err)
					}

					if err := i18n.ReloadLangs(); err != nil {
						log.Error("Conf Reload: ", err)
					}

					reloadConfig()
					log.Info("Config Reloaded")

				case ".json":
					if checkEventTime(event.Name) {
						continue
					}
					if event.Name == CompressConfPath {
						settingCompress()
						log.Info("Beego Compress Reloaded")
					}
				}
			}
		}
	}()

	if err := watcher.WatchFlags("conf", fsnotify.FSN_MODIFY); err != nil {
		log.Error(err)
	}

	if err := watcher.WatchFlags("conf/global", fsnotify.FSN_MODIFY); err != nil {
		log.Error(err)
	}
}

// checkEventTime returns true if FileModTime does not change.
func checkEventTime(name string) bool {
	mt := getFileModTime(name)
	if eventTime[name] == mt {
		return true
	}

	eventTime[name] = mt
	return false
}

// getFileModTime retuens unix timestamp of `os.File.ModTime` by given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		log.Error("Fail to open file[ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Error("Fail to get file information[ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

func IsMatchHost(uri string) bool {
	if len(uri) == 0 {
		return false
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	if u.Host != AppHost {
		return false
	}

	return true
}

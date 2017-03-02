captcha [![Build Status](https://drone.io/github.com/tango-contrib/captcha/status.png)](https://drone.io/github.com/tango-contrib/captcha/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/captcha)](http://gocover.io/github.com/tango-contrib/captcha)
====

Middleware captcha a middleware that provides captcha service for [Tango](https://github.com/lunny/tango).

[API Reference](https://gowalker.org/github.com/tango-contrib/captcha)

### Installation

	go get github.com/tango-contrib/captcha
	
## Usage

```go
// main.go
import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/cache"
	"github.com/tango-contrib/captcha"
)

type CaptchaAction struct {
	captcha.Captcha
	renders.Renderer
}

func (c *CaptchaAction) Get() {
	c.Render("captcha.html", renders.T{
		"captcha": c.CreateHtml(),
	})
}

func (c *CaptchaAction) Post() string {
	if c.Verify() {
		return "true"
	}
	return "false"
}

func main() {
  	t := tango.Classic()
	t.Use(captcha.New())
	t.Any("/", new(CaptchaAction))
	t.Run()
}
```

```html
<!-- templates/captcha.tmpl -->
{{.captcha}}
```

## Options

`captcha.Captchaer` comes with a variety of configuration options:

```go
// ...
t.Use(captcha.New(captcha.Options{
	URLPrefix:			"/captcha/", 	// URL prefix of getting captcha pictures.
	FieldIdName:		"captcha_id", 	// Hidden input element ID.
	FieldCaptchaName:	"captcha", 		// User input value element name in request form.
	ChallengeNums:		6, 				// Challenge number.
	Width:				240,			// Captcha image width.
	Height:				80,				// Captcha image height.
	Expiration:			600, 			// Captcha expiration time in seconds.
	CachePrefix:		"captcha_", 	// Cache key prefix captcha characters.
}, cache))
// ...
```

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.
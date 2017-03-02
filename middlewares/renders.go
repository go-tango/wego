package middlewares

import (
	"html/template"

	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/setting"

	"github.com/tango-contrib/renders"
)

var (
	Renders *renders.Renders
)

func Init() {
	Renders = renders.New(renders.Options{
		Directory: setting.TemplatesPath,
		Funcs:     mergeFuncMap(utils.FuncMap(), setting.Funcs),
		Vars: renders.T{
			"AppName":       setting.AppName,
			"AppVer":        setting.AppVer,
			"AppUrl":        setting.AppUrl,
			"AppLogo":       setting.AppLogo,
			"AvatarURL":     setting.AvatarURL,
			"IsProMode":     setting.IsProMode,
			"SearchEnabled": setting.SearchEnabled,
		},
	})
}

func mergeFuncMap(funcs ...template.FuncMap) template.FuncMap {
	var ret = make(template.FuncMap)
	for _, fs := range funcs {
		for k, f := range fs {
			ret[k] = f
		}
	}
	return ret
}

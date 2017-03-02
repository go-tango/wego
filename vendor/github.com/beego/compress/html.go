package compress

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Group struct {
	DistFile    string
	SourceFiles []string
	SkipFiles   []string
}

type compress struct {
	StaticURL string
	SrcPath   string
	DistPath  string
	SrcURL    string
	DistURL   string
	Groups    map[string]Group
	IsProMode bool
	caches    map[string]template.HTML
}

func (c *compress) SetProMode(isPro bool) {
	c.IsProMode = isPro
}

func (c *compress) SetStaticURL(url string) {
	c.StaticURL = url
}

func errHtml(err string, args ...interface{}) string {
	err = fmt.Sprintf("Beego Compress: "+err, args...)
	fmt.Fprintln(os.Stderr, err)
	return "<!-- " + err + " -->"
}

type compressJs struct {
	compress
}

func (c *compressJs) CompressJs(name string) template.HTML {
	return generateHTML(name, c.compress, JsTagTemplate)
}

type compressCss struct {
	compress
}

func (c *compressCss) CompressCss(name string) template.HTML {
	return generateHTML(name, c.compress, CssTagTemplate)
}

func generateHTML(name string, c compress, t *template.Template) template.HTML {
	if group, ok := c.Groups[name]; ok {
		if c.IsProMode {

			if c.caches == nil {
				c.caches = make(map[string]template.HTML, len(c.Groups))
			}

			if scripts, ok := c.caches[name]; ok {
				return scripts
			}

			scripts := fmt.Sprintf("<script>/* Beego Compress Powered */</script>\n\t")

			filePath := filepath.Join(c.DistPath, group.DistFile)
			if info, err := os.Stat(filePath); err == nil {
				URL := c.StaticURL + path.Join(c.DistURL, group.DistFile) + "?ver=" + fmt.Sprint(info.ModTime().Unix())

				if res, err := parseTmpl(t, map[string]string{"URL": URL}); err != nil {
					errHtml("tempalte execute error: %s", err)

				} else {
					scripts += res
				}

			} else {
				errHtml("load file `%s` for path `%s` error: %s", group.DistFile, filePath, err.Error())
			}

			if len(scripts) > 0 {
				res := template.HTML(scripts + "\n")
				c.caches[name] = res
				return res
			}
		} else {
			scripts := make([]string, 0, len(group.SourceFiles)+2)

			scripts = append(scripts, fmt.Sprintf("<script>/* Beego Compress group `%s` begin */</script>", name))

			for _, file := range group.SourceFiles {
				filePath := filepath.Join(c.SrcPath, file)

				if info, err := os.Stat(filePath); err == nil {
					URL := c.StaticURL + path.Join(c.SrcPath, file) + "?ver=" + fmt.Sprint(info.ModTime().Unix())

					if res, err := parseTmpl(t, map[string]string{"URL": URL}); err != nil {
						scripts = append(scripts, errHtml("tempalte execute error: %s", err))

					} else {
						scripts = append(scripts, res)
					}

				} else {
					scripts = append(scripts, errHtml("load file `%s` for path `%s` error: %s", file, filePath, err.Error()))
				}
			}

			scripts = append(scripts, fmt.Sprintf("<script>/* end */</script>"))

			return template.HTML(strings.Join(scripts, "\n\t"))
		}
	} else {
		return template.HTML(errHtml("not found compress group `%s`", name))
	}

	return ""
}

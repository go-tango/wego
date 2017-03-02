package compress

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func compressJsFiles(s *Settings, force, skip, verbose bool) {
	compressFiles(&s.Js.(*compressJs).compress, force, skip, verbose, JsFilters)
}

func compressCssFiles(s *Settings, force, skip, verbose bool) {
	compressFiles(&s.Css.(*compressCss).compress, force, skip, verbose, CssFilters)
}

func logError(err string, args ...interface{}) {
	err = fmt.Sprintf(err, args...)
	fmt.Fprintln(os.Stderr, err)
}

func logInfo(info string, args ...interface{}) {
	info = fmt.Sprintf(info, args...)
	fmt.Fprintln(os.Stdout, info)
}

func compressFiles(c *compress, force, skip, verbose bool, filters []Filter) {
	os.MkdirAll(TmpPath, 0755)

	for name, group := range c.Groups {

		hasError := false
		hasModified := false
		sources := make([]string, 0, len(group.SourceFiles))

		if verbose {
			logInfo("Group '%s'", name)
			logInfo("--------------------------")
		}

		skips := make(map[string]bool, len(group.SkipFiles))
		for _, file := range group.SkipFiles {
			skips[file] = true
		}

		for _, file := range group.SourceFiles {

			modified := false

			var cacheTime *time.Time
			cacheFile := filepath.Join(TmpPath, c.SrcPath, file)
			if info, err := os.Stat(cacheFile); err == nil {
				// get cached file modtime
				t := info.ModTime()
				cacheTime = &t
			}

			sourceFile := filepath.Join(c.SrcPath, file)
			if info, err := os.Stat(sourceFile); err == nil {
				if cacheTime != nil {
					if info.ModTime().Unix() > cacheTime.Unix() {
						// file modified
						modified = true
					}
				} else {
					modified = true
				}
			} else {
				logError("source file %s load error: %s", sourceFile, err.Error())
				hasError = true
				continue
			}

			if skip || modified {
				buf := bytes.NewBufferString("")
				// load content from file
				if f, err := os.Open(sourceFile); err == nil {
					buf.ReadFrom(f)
					f.Close()
				} else {
					logError("source file %s load error: %s", sourceFile, err.Error())
					hasError = true
					continue
				}

				source := buf.String()
				if verbose {
					fmt.Fprintf(os.Stdout, "compress file %s ... ", sourceFile)
				}
				if skips[file] {
					if verbose {
						fmt.Fprintf(os.Stdout, "skiped ")
					}
				} else {
					for _, filter := range filters {
						// compress content
						source = filter(source)
					}
				}
				sources = append(sources, source)

				var writeErr error
				dir, _ := filepath.Split(cacheFile)
				if writeErr = os.MkdirAll(dir, 0755); writeErr == nil {
					if f, err := os.OpenFile(cacheFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
						if _, err := f.WriteString(source); err == nil {
							hasModified = true
							if verbose {
								logInfo("saved")
							}
						} else {
							writeErr = err
						}
						f.Close()
					} else {
						writeErr = err
					}
				}

				if writeErr != nil {
					logError("write error: %s", writeErr.Error())
					hasError = true
				}

			} else {
				buf := bytes.NewBufferString("")
				// load content from file
				if f, err := os.Open(cacheFile); err == nil {
					buf.ReadFrom(f)
				} else {
					logError("cache file %s load error: %s", cacheFile, err.Error())
					hasError = true
					continue
				}

				if verbose {
					logInfo("use cache file %s", cacheFile)
				}
				sources = append(sources, buf.String())
			}
		}

		if !hasError {
			if hasModified || force {
				distFile := filepath.Join(c.DistPath, group.DistFile)
				var writeErr error
				dir, _ := filepath.Split(distFile)
				if writeErr = os.MkdirAll(dir, 0755); writeErr == nil {
					if f, err := os.OpenFile(distFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
						if _, err := f.WriteString(strings.Join(sources, "\n\n")); err == nil {
							if verbose {
								logInfo("compressed file %s ... saved", distFile)
							}
						} else {
							writeErr = err
						}
						f.Close()

					} else {
						writeErr = err
					}
				}

				if writeErr != nil {
					logError("compressed file %s write error: %s", distFile, writeErr.Error())
					hasError = true
				}
			} else {
				if verbose {
					logInfo("not modified")
				}
			}
		}

		if verbose {
			logInfo("")
		}
	}
}

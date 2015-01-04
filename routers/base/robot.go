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

package base

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/go-tango/wego/setting"
)

var robotTxt string

const robotTpl = `{{$disallow := .Disallow}}{{range .Uas}}User-Agent: {{.}}
Disallow: {{$disallow}}

{{end}}User-Agent: *
Disallow: /
`

// RobotRouter implemented global settings for all other routers.
type RobotRouter struct {
}

// Get implemented Prepare method for RobotRouter.
func (this *RobotRouter) Get() string {
	if len(robotTxt) > 0 {
		return robotTxt
	}

	// Generate "robot.txt".
	t := template.New("robotTpl")
	t.Parse(robotTpl)
	uas := strings.Split(setting.Cfg.MustValue("robot", "uas"), "|")

	buf := new(bytes.Buffer)
	t.Execute(buf, map[string]interface{}{
		"Uas": uas,
		"Disallow": setting.Cfg.MustValue("robot", "disallow"),
	})
	return buf.String()
}

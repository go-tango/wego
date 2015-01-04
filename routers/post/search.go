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

package post

import (
	"fmt"
	"github.com/go-tango/wego/routers/base"
)

type SearchRouter struct {
	base.BaseRouter
}

func (this *SearchRouter) Get() {
	site := "site:golanghome.com"
	q := this.GetString("q")
	this.Redirect(fmt.Sprintf("http://cn.bing.com/search?q=%s %s", site, q), 302)
}

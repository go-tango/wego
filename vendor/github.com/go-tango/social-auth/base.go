// Copyright 2014 beego authors
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
//
// Maintain by https://github.com/slene

package social

import (
	"fmt"
)

var providers = make(map[SocialType]Provider)
var providersByPath = make(map[string]Provider)

// Register the social provider
func RegisterProvider(prov Provider) error {
	typ := prov.GetType()
	if !typ.Available() {
		return fmt.Errorf("Unknown social type `%d`", typ)
	}
	path := prov.GetPath()
	if providersByPath[path] != nil {
		return fmt.Errorf("path `%s` is already in used", path)
	}
	providers[typ] = prov
	providersByPath[path] = prov
	return nil
}

// Get provider by SocialType
func GetProviderByType(typ SocialType) (Provider, bool) {
	if p, ok := providers[typ]; ok {
		return p, true
	} else {
		return nil, false
	}
}

// Get provider by path name
func GetProviderByPath(path string) (Provider, bool) {
	if p, ok := providersByPath[path]; ok {
		return p, true
	} else {
		return nil, false
	}
}

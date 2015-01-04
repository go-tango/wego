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

package auth

import (
	"net/http"
	"bytes"
	"fmt"
	"strings"
	"encoding/base64"
	"crypto/sha1"
	"strconv"
	"time"
	"crypto/hmac"

	"github.com/astaxie/beego/orm"

	"github.com/go-tango/wego/modules/models"
)

func GetSecureCookie(req *http.Request, Secret, key string) (string, bool) {
	val := GetCookie(req, key)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// Set Secure cookie for response.
func SetSecureCookie(resp http.ResponseWriter, Secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	SetCookie(resp, name, cookie, others...)
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

var cookieValueSanitizer = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func sanitizeValue(v string) string {
	return cookieValueSanitizer.Replace(v)
}

func SetCookie(resp http.ResponseWriter, name string, value string, others ...interface{}) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", sanitizeName(name), sanitizeValue(value))
	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int64:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int32:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		}
	}

	// the settings below
	// Path, Domain, Secure, HttpOnly
	// can use nil skip set

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Path=%s", sanitizeValue(v))
		}
	} else {
		fmt.Fprintf(&b, "; Path=%s", "/")
	}

	// default empty
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Domain=%s", sanitizeValue(v))
		}
	}

	// default empty
	if len(others) > 3 {
		var secure bool
		switch v := others[3].(type) {
		case bool:
			secure = v
		default:
			if others[3] != nil {
				secure = true
			}
		}
		if secure {
			fmt.Fprintf(&b, "; Secure")
		}
	}

	// default false. for session cookie default true
	httponly := false
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			// HttpOnly = true
			httponly = true
		}
	}

	if httponly {
		fmt.Fprintf(&b, "; HttpOnly")
	}

	resp.Header().Add("Set-Cookie", b.String())
}

func GetCookie(req *http.Request, key string) string {
	ck, err := req.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

func UserFollow(user *models.User, theUser *models.User) {
	if theUser.Read() == nil {
		var mutual bool
		tFollow := models.Follow{User: theUser, FollowUser: user}
		if err := tFollow.Read("User", "FollowUser"); err == nil {
			mutual = true
		}

		follow := models.Follow{User: user, FollowUser: theUser, Mutual: mutual}
		if err := follow.Insert(); err == nil && mutual {
			tFollow.Mutual = mutual
			tFollow.Update("Mutual")
		}

		if nums, err := user.FollowingUsers().Count(); err == nil {
			user.Following = int(nums)
			user.Update("Following")
		}

		if nums, err := theUser.FollowerUsers().Count(); err == nil {
			theUser.Followers = int(nums)
			theUser.Update("Followers")
		}
	}
}

func UserUnFollow(user *models.User, theUser *models.User) {
	num, _ := user.FollowingUsers().Filter("FollowUser", theUser.Id).Delete()
	if num > 0 {
		theUser.FollowingUsers().Filter("FollowUser", user.Id).Update(orm.Params{
			"Mutual": false,
		})

		if nums, err := user.FollowingUsers().Count(); err == nil {
			user.Following = int(nums)
			user.Update("Following")
		}

		if nums, err := theUser.FollowerUsers().Count(); err == nil {
			theUser.Followers = int(nums)
			theUser.Update("Followers")
		}
	}
}

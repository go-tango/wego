package models

import "github.com/go-tango/wego/modules/utils"

// return a user salt token
func GetUserSalt() string {
	return utils.GetRandomString(10)
}

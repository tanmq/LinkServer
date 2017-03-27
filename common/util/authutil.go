package util

import (
	"linkServer/packet"
)

//Auth is to auth the user is valid or not.
func Auth(loginInfo *packet.LoginInfo) bool {
	return true
}

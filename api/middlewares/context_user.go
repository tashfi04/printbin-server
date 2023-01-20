package middlewares

import "net/http"

// UserInfoCxt prefix
const UserInfoCxt = "userinfo_"

// UserInfo stores Token, UserID and CityID from the request header
type UserInfo struct {
	UserID uint
	Role   uint
}

// GetUserInfo returns UserInfo after adding UserInfoCtx as context
func GetUserInfo(r *http.Request) (*UserInfo, error) {

	ui := r.Context().Value(UserInfoCxt)

	return ui.(*UserInfo), nil
}

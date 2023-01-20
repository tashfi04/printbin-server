package api

import (
	"github.com/tashfi04/printbin-server/api/middlewares"
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/dtos"
	"github.com/tashfi04/printbin-server/repos"
	"github.com/tashfi04/printbin-server/session"
	"github.com/tashfi04/printbin-server/utils"
	"gorm.io/gorm"
	"net/http"
)

func isAuthenticated(w http.ResponseWriter, r *http.Request) {

	userInfo, err := middlewares.GetUserInfo(r)
	if err != nil {
		utils.Logger().Errorln("Missing user context in header: ", err)
		rndr.JSON(w, http.StatusUnprocessableEntity, utils.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Missing user context in header",
			Error:      "Missing user context in header",
		})
		return
	}

	dbResponse, err := repos.UserRepo().GetMinimalUserInfo(conn.DB(), userInfo.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger().Errorln(err)
			rndr.JSON(w, http.StatusNotFound, utils.Response{
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
				Error:      err.Error(),
			})
			return
		}
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database Query failed",
			Error:      err.Error(),
		})
		return
	}

	rndr.JSON(w, http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Data:       dbResponse,
		Message:    "Authorized user",
	})
	return
}

func login(w http.ResponseWriter, r *http.Request) {

	username := r.Header.Get("Username")
	if username == "" {
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Username missing in header",
		})
		return
	}

	password := r.Header.Get("Password")
	if password == "" {
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Password missing in header",
		})
		return
	}

	userInfo, err := repos.UserRepo().GetUserByUsername(conn.DB(), username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rndr.JSON(w, http.StatusNotFound, utils.Response{
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
				Error:      err.Error(),
			})
			return
		}
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database query failed",
			Error:      err.Error(),
		})
		return
	}

	userPassword, err := utils.Decrypt(userInfo.Password)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Server error",
			Error:      err.Error(),
		})
		return
	}

	if password != userPassword {
		utils.Logger().Errorln("invalid password")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid password",
			Error:      "invalid password",
		})
		return
	}

	userSession, err := session.CookieStore().Get(r, "auth_token")
	userSession.Values["authenticated"] = true
	userSession.Values["id"] = userInfo.ID
	userSession.Values["role"] = userInfo.Role
	if err = userSession.Save(r, w); err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create user session",
			Error:      err.Error(),
		})
		return
	}

	resp := dtos.MinimalUserInfo{
		Username: userInfo.Username,
		Role:     userInfo.Role,
	}

	rndr.JSON(w, http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Data:       resp,
		Message:    "User session created successfully",
	})
	return
}

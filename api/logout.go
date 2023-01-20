package api

import (
	"github.com/tashfi04/printbin-server/session"
	"github.com/tashfi04/printbin-server/utils"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {

	userSession, _ := session.CookieStore().Get(r, "auth_token")

	delete(userSession.Values, "authenticated")
	delete(userSession.Values, "id")
	delete(userSession.Values, "role")
	userSession.Options.MaxAge = -1
	err := userSession.Save(r, w)
	if err != nil {
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error while logging out",
			Error:      err.Error(),
		})
		return
	}

	rndr.JSON(w, http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Message:    "Successfully logged out",
	})
	return
}

package middlewares

import (
	"github.com/tashfi04/printbin-server/utils"
	"net/http"
)

func AuthenticateAdmin(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userInfo, err := GetUserInfo(r)
		if err != nil {
			utils.Logger().Errorln("Missing user context in header: ", err)
			rndr.JSON(w, http.StatusUnprocessableEntity, utils.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Missing user context in header",
				Error:      "Missing user context in header",
			})
			return
		}

		if userInfo.Role < 1 {
			utils.Logger().Infoln("Unauthorized access prohibited")
			rndr.JSON(w, http.StatusUnauthorized, utils.Response{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized access prohibited",
				Error:      "Unauthorized access prohibited",
			})
			return
		}

		next.ServeHTTP(w, r)
	})

}

package middlewares

import (
	"context"
	"github.com/tashfi04/printbin-server/session"
	"github.com/tashfi04/printbin-server/utils"
	"github.com/thedevsaddam/renderer"
	"net/http"
)

var (
	rndr *renderer.Render
)

func init() {
	rndr = renderer.New()
}

func AuthenticateClient(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userSession, _ := session.CookieStore().Get(r, "auth_token")

		// Check if user is authenticated
		if auth, ok := userSession.Values["authenticated"].(bool); !ok || !auth {
			rndr.JSON(w, http.StatusUnauthorized, utils.Response{
				StatusCode: http.StatusUnauthorized,
				Message:    "Must be logged in to view this page",
				Error:      "Unauthorized access prohibited",
			})
			return
		}

		userID, ok := userSession.Values["id"].(uint)
		if !ok {
			rndr.JSON(w, http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch user id from session",
				Error:      "Failed to fetch user id from session",
			})
			return
		}

		role, ok := userSession.Values["role"].(uint)
		if !ok {
			rndr.JSON(w, http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch admin status from session",
				Error:      "Failed to fetch admin status from session",
			})
			return
		}

		userInfo := &UserInfo{
			UserID: userID,
			Role:   role,
		}

		ctx := context.WithValue(r.Context(), UserInfoCxt, userInfo)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

package session

import (
	"github.com/gorilla/sessions"
	"github.com/tashfi04/printbin-server/config"
)

var store *sessions.CookieStore

func NewCookieStore() error {
	cfg := config.Session()
	store = sessions.NewCookieStore([]byte(cfg.SessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 1,
		HttpOnly: true,
		Secure:   false,
		//Domain:   config.App().HomepageUrl,
	}
	return nil
}

func CookieStore() *sessions.CookieStore {
	return store
}

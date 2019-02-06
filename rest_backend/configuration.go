package rest_backend

import "github.com/gorilla/securecookie"

type Configuration struct {
	SecureCookies *securecookie.SecureCookie
	HttpsCookies bool
}

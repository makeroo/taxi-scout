package rest_backend


type CookieManager interface {
	Encode(name string, value interface{}) (string, error)
	Decode(name, value string, dst interface{}) error
}


type Configuration struct {
	SecureCookies CookieManager
	HttpsCookies bool
}

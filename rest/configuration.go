package rest

// CookieManager abstracts cookie manipulation functions such encryption and decryption.
type CookieManager interface {
	Encode(name string, value interface{}) (string, error)
	Decode(name, value string, dst interface{}) error
}

// Configuration provides REST backend configuration.
type Configuration struct {
	SecureCookies CookieManager
	HTTPSCookies  bool
}

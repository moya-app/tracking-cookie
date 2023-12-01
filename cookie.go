package tracking_cookie

import (
    "reflect"
	"context"
	"net/http"
    "time"
    "math/rand"
    "encoding/base64"
    "os"
    "fmt"
)

const defaultCookieName = "reqId"

// Config the plugin configuration.
type Config struct {
    Domain string `json:"domain,omitempty" yaml:"domain,omitempty" toml:"domain,omitempty"`
    Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
    Expires int `json:"expires,omitempty" yaml:"expires,omitempty" toml:"expires,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
        Name: defaultCookieName,
        // 100 year expiry date
        Expires: 100 * 365 * 24 * 60 * 60,
	}
}

type UserCookies struct {
	next     http.Handler
	name     string
    domain   string
    cookieName string
    cookieExpires int
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &UserCookies{
		next:     next,
		name:     name,
        domain:   config.Domain,
        cookieName: config.Name,
        cookieExpires: config.Expires,
	}, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)

	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomString)
}

func (a *UserCookies) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    cookie, err := req.Cookie(a.cookieName)
    if err != nil {
        cookie = &http.Cookie{
            Name: a.cookieName,
            Value: generateRandomString(20),
            Domain: a.domain,

            // Send on all paths
            Path: "/",

            // Set expiry a long time in the future
            MaxAge: a.cookieExpires,
            Expires: time.Now().Add(time.Duration(a.cookieExpires) * time.Second),

            // Don't allow JS access to this
            HttpOnly: true,
        }

        // Generate a a new cookie for this user - havn't seen them before
        http.SetCookie(rw, cookie)
    }

    // TODO: Export cookie.Value via a custom accesslog parameter
	a.next.ServeHTTP(rw, req)
}

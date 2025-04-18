package httpclient

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/valyala/fasthttp"

	"github.com/oidc-mytoken/utils/context"
)

var client *resty.Client

func init() {
	client = resty.New()
	client.SetCookieJar(nil)
	// client.SetDisableWarn(true)
	client.SetRetryCount(2)
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))
	client.SetTimeout(20 * time.Second)
	context.SetClient(client.GetClient())
}

// Init initializes the http client
func Init(hostURL, userAgent string) {
	if hostURL != "" {
		client.SetBaseURL(hostURL)
	}
	client.SetHeader(fasthttp.HeaderUserAgent, userAgent)
}

// Do returns the client, so it can be used to do requests
func Do() *resty.Client {
	return client
}

// EnableDebug enables debug on the client
func EnableDebug() {
	client.SetDebug(true)
}

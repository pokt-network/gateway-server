package relayer

import (
	"github.com/valyala/fasthttp"
	"time"
)

type httpRequester interface {
	DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error
}

// fastHttpRequester: used to mock fast http for testing
type fastHttpRequester struct{}

func (receiver fastHttpRequester) DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error {
	return fasthttp.DoTimeout(req, resp, timeout)
}

package middleware

import (
	"github.com/valyala/fasthttp"
	"pokt_gateway_server/cmd/gateway_server/internal/common"
	config2 "pokt_gateway_server/internal/global_config"
)

func retrieveAPIKey(ctx *fasthttp.RequestCtx) string {
	auth := ctx.Request.Header.Peek("x-api-key")
	if auth == nil {
		return ""
	}
	return string(auth)
}

// BasicAuth is the basic auth handler
func XAPIKeyAuth(h fasthttp.RequestHandler, provider config2.SecretProvider) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Get the Basic Authentication credentials
		xAPIKey := retrieveAPIKey(ctx)

		if xAPIKey != "" && xAPIKey == provider.GetAPIKey() {
			h(ctx)
			return
		}
		// Request Basic Authentication otherwise
		common.JSONError(ctx, "Unauthorized, invalid x-api-key header", fasthttp.StatusUnauthorized)
	}
}

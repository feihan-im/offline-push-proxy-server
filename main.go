package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/feihan-im/offline-push-proxy-server/handler"
)

func main() {
	port := flag.Int("port", 21001, "port to listen")
	apiURL := flag.String("apiURL", "https://push-api.feihanim.cn", "API URL to forward requests")
	flag.Parse()

	h := server.Default(server.WithHostPorts(fmt.Sprintf(":%d", *port)))

	h.GET("/v1/token", func(ctx context.Context, c *app.RequestContext) {
		handler.TokenHandler(ctx, c, *apiURL)
	})
	h.GET("/v1/meta", func(ctx context.Context, c *app.RequestContext) {
		handler.MetaHandler(ctx, c, *apiURL)
	})
	h.POST("/v1/push", func(ctx context.Context, c *app.RequestContext) {
		handler.PushHandler(ctx, c, *apiURL)
	})

	h.Spin()
}

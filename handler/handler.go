package handler

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/network/standard"
)

var globalClient *client.Client

func init() {
	// 全局初始化 Hertz 客户端
	client, err := client.NewClient(
		client.WithDialer(standard.NewDialer()),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}
	globalClient = client
}

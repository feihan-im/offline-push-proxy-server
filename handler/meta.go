package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func MetaHandler(ctx context.Context, c *app.RequestContext, apiURL string) {
	statusCode, body, err := globalClient.Get(ctx, nil, apiURL+"/v1/meta")
	if err != nil {
		hlog.Error("Failed to forward request", err)
		c.JSON(consts.StatusInternalServerError, utils.H{"code": 401, "msg": "failed to forward request"})
		return
	}

	c.Data(statusCode, "application/json", body)
}

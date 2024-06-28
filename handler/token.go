package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type TokenGetReq struct {
	OrgCode   string `query:"orgCode"`
	Signature string `query:"signature"`
	Timestamp int64  `query:"timestamp"`
	Nonce     int64  `query:"nonce"`
}

func TokenHandler(ctx context.Context, c *app.RequestContext, apiURL string) {
	var req TokenGetReq
	if err := c.BindQuery(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": err.Error()})
		return
	}

	if len(req.OrgCode) > 32 {
		c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "orgCode too long"})
		return
	}

	if len(req.Signature) > 100 {
		c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "signature too long"})
		return
	}

	statusCode, body, err := globalClient.Get(ctx, nil, apiURL+"/v1/token?"+c.QueryArgs().String())
	if err != nil {
		hlog.Error("Failed to forward request", err)
		c.JSON(consts.StatusInternalServerError, utils.H{"code": 401, "msg": "failed to forward request"})
		return
	}

	c.Data(statusCode, "application/json", body)
}

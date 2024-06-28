package handler

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type PushReq struct {
	Devices []*Device `json:"devices"`
	Type    string    `json:"type"`
	Msg     *Msg      `json:"msg,omitempty"`
}

type Device struct {
	Platform    string `json:"platform"`
	DeviceToken string `json:"deviceToken"`
	Development bool   `json:"development,omitempty"`
}

type Msg struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	Group       string `json:"group,omitempty"`
	GroupPrefix string `json:"groupPrefix,omitempty"`
	ClickType   string `json:"clickType"`
	Intent      string `json:"intent,omitempty"`
	Badge       int64  `json:"badge,omitempty"`
	Sound       bool   `json:"sound,omitempty"`
}

func PushHandler(ctx context.Context, c *app.RequestContext, apiURL string) {
	var req PushReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": err.Error()})
		return
	}

	if len(req.Type) > 10 {
		c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "type too long"})
		return
	}

	for _, device := range req.Devices {
		if len(device.Platform) > 20 {
			c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "platform too long"})
			return
		}
		if len(device.DeviceToken) > 200 {
			c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "deviceToken too long"})
			return
		}
	}

	if req.Msg != nil {
		if len(req.Msg.Title) > 150 {
			c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "title too long"})
			return
		}
		if len(req.Msg.Body) > 150 {
			c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "body too long"})
			return
		}
		if len(req.Msg.Group) > 200 {
			c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "group too long"})
			return
		}
		if len(req.Msg.GroupPrefix) > 150 {
			c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "groupPrefix too long"})
			return
		}
		if req.Msg.ClickType != "intent" {
			c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "invalid clickType"})
			return
		}
		if len(req.Msg.Intent) > 200 {
			c.JSON(consts.StatusBadRequest, utils.H{"code": 401, "msg": "intent too long"})
			return
		}
	}

	payload, err := sonic.Marshal(req)
	if err != nil {
		hlog.Error("Failed to marshal request body", err)
		c.JSON(consts.StatusInternalServerError, utils.H{"code": 401, "msg": "failed to marshal request body"})
		return
	}

	httpReq, httpRes := protocol.AcquireRequest(), protocol.AcquireResponse()
	httpReq.SetRequestURI(apiURL + "/v1/push")
	httpReq.SetMethod("POST")
	httpReq.SetHeader(consts.HeaderAuthorization, c.Request.Header.Get(consts.HeaderAuthorization))
	httpReq.SetBody(payload)

	if err := globalClient.Do(ctx, httpReq, httpRes); err != nil {
		hlog.Error("Failed to forward request", err)
		c.JSON(consts.StatusInternalServerError, utils.H{"code": 401, "msg": "failed to forward request"})
	}

	c.Data(httpRes.StatusCode(), "application/json", httpRes.Body())
}

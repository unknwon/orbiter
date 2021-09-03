// Copyright 2016 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package routers

import (
	"bytes"
	"encoding/json"

	"unknwon.dev/orbiter/internal/context"
	"unknwon.dev/orbiter/internal/models"
	"unknwon.dev/orbiter/internal/models/errors"
)

func Webhooks(ctx *context.Context) {
	ctx.Data["Title"] = "Webhooks"
	ctx.Data["PageIsWebhook"] = true

	webhooks, err := models.QueryWebhooks(models.QueryWebhookOptions{
		Limit: 50,
		Order: "created desc",
	})
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Data["Webhooks"] = webhooks

	ctx.HTML(200, "webhook/list")
}

func ViewWebhook(ctx *context.Context) {
	ctx.Data["Title"] = "View Webhook"
	ctx.Data["PageIsWebhook"] = true
	ctx.Data["RequireHighlightJS"] = true

	webhook, err := models.GetWebhookByID(ctx.ParamsInt64(":id"))
	if err != nil {
		if errors.IsWebhookNotFound(err) {
			ctx.Handle(404, "GetWebhookByID", nil)
		} else {
			ctx.Handle(500, "GetWebhookByID", err)
		}
		return
	}
	ctx.Data["Webhook"] = webhook

	// Prettify JSON in case it is not.
	buf := new(bytes.Buffer)
	if err = json.Indent(buf, []byte(webhook.Payload), "", "  "); err != nil {
		ctx.Handle(500, "json.Indent", err)
		return
	}
	webhook.Payload = buf.String()

	ctx.HTML(200, "webhook/view")
}

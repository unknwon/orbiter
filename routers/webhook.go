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
	"github.com/Unknwon/orbiter/models"
	"github.com/Unknwon/orbiter/modules/context"
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

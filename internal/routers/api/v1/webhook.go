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

package v1

import (
	"unknwon.dev/orbiter/internal/models"
)

func ListWebhooks(ctx *Context) {
	webhooks, err := models.QueryWebhooks(models.QueryWebhookOptions{
		CollectorID: ctx.QueryInt64("collector_id"),
		Owner:       ctx.Query("owner"),
		RepoName:    ctx.Query("repo_name"),
		EventType:   ctx.Query("event_type"),
		Sender:      ctx.Query("sender"),
		After:       ctx.QueryInt64("after"),
		Limit:       ctx.QueryInt64("limit"),
	})
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.JSON(200, webhooks)
}

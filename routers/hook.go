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
	"github.com/Unknwon/orbiter/modules/webhook"
)

func Hook(ctx *context.Context) {
	collector, err := models.GetCollectorBySecret(ctx.Query("secret"))
	if err != nil {
		if models.IsErrCollectorNotFound(err) {
			ctx.Error(403)
		} else {
			ctx.Error(500, err.Error())
		}
		return
	}

	payload, err := ctx.Req.Body().Bytes()
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	// NOTE: Currently only support GitHub
	event, err := webhook.ParseGitHubEvent(payload)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	if err = models.NewWebhook(&models.Webhook{
		CollectorID: collector.ID,
		Owner:       event.Repository.Owner.Login,
		RepoName:    event.Repository.Name,
		EventType:   ctx.Req.Header.Get("X-GitHub-Event"),
		Sender:      event.Sender.Login,
		Payload:     string(payload),
	}); err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Status(202)
}

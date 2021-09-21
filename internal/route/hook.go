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

package route

import (
	"net/http"

	"github.com/flamego/flamego"

	"unknwon.dev/orbiter/internal/db"
	"unknwon.dev/orbiter/internal/db/errors"
	"unknwon.dev/orbiter/internal/tool"
	"unknwon.dev/orbiter/internal/webhook"
)

// POST /hook
func Hook(c flamego.Context) {
	collector, err := db.GetCollectorBySecret(c.Query("secret"))
	if err != nil {
		if errors.IsCollectorNotFound(err) {
			c.ResponseWriter().WriteHeader(http.StatusForbidden)
		} else {
			http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
		return
	}

	payload, err := c.Request().Body().Bytes()
	if err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE: Currently only support GitHub
	event, err := webhook.ParseGitHubEvent(payload)
	if err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}

	if err = db.NewWebhook(&db.Webhook{
		CollectorID: collector.ID,
		Owner:       tool.FirstNonEmptyString(event.Repository.Owner.Login, event.Repository.Owner.Name),
		RepoName:    event.Repository.Name,
		EventType:   c.Request().Header.Get("X-GitHub-Event"),
		Sender:      event.Sender.Login,
		Payload:     string(payload),
	}); err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.ResponseWriter().WriteHeader(http.StatusAccepted)
}

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
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"unknwon.dev/orbiter/internal/db"
	"unknwon.dev/orbiter/internal/db/errors"
)

// GET /webhooks
func Webhooks(w http.ResponseWriter, t template.Template, data template.Data) {
	data["Title"] = "Webhooks"
	data["PageIsWebhook"] = true

	webhooks, err := db.QueryWebhooks(
		db.QueryWebhookOptions{
			Limit: 50,
			Order: "created desc",
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Webhooks"] = webhooks

	t.HTML(http.StatusOK, "webhook/list")
}

// GET /webhooks/{id}
func ViewWebhook(c flamego.Context, t template.Template, data template.Data) {
	data["Title"] = "View Webhook"
	data["PageIsWebhook"] = true
	data["RequireHighlightJS"] = true

	webhook, err := db.GetWebhookByID(c.ParamInt64("id"))
	if err != nil {
		if errors.IsWebhookNotFound(err) {
			http.NotFound(c.ResponseWriter(), c.Request().Request)
		} else {
			http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
		return
	}
	data["Webhook"] = webhook

	// Prettify JSON in case it is not.
	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(webhook.Payload), "", "  ")
	if err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}
	webhook.Payload = buf.String()

	t.HTML(http.StatusOK, "webhook/view")
}

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
	"encoding/json"
	"net/http"

	"unknwon.dev/orbiter/internal/db"
)

func ListWebhooks(c *Context) {
	webhooks, err := db.QueryWebhooks(db.QueryWebhookOptions{
		CollectorID: c.QueryInt64("collector_id"),
		Owner:       c.Query("owner"),
		RepoName:    c.Query("repo_name"),
		EventType:   c.Query("event_type"),
		Sender:      c.Query("sender"),
		After:       c.QueryInt64("after"),
		Limit:       c.QueryInt64("limit"),
	})
	if err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	_ = json.NewEncoder(c.ResponseWriter()).Encode(webhooks)
}

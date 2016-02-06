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

package models

import (
	"time"
)

// Webhook represents a history record of webhook.
type Webhook struct {
	ID          int64
	CollectorID int64 `sql:"index"`
	Owner       string
	RepoName    string
	EventType   string
	Sender      string
	Payload     string `sql:"type:text"`
	Created     int64
}

func NewWebhook(webhook *Webhook) error {
	webhook.Created = time.Now().UTC().UnixNano()
	return x.Create(webhook).Error
}

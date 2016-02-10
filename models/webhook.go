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

	log "github.com/Sirupsen/logrus"

	"github.com/Unknwon/orbiter/models/errors"
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

func (w *Webhook) CreatedTime() time.Time {
	return time.Unix(0, w.Created)
}

func NewWebhook(webhook *Webhook) error {
	webhook.Created = time.Now().UTC().UnixNano()
	return x.Create(webhook).Error
}

func GetWebhookByID(id int64) (*Webhook, error) {
	webhook := new(Webhook)
	err := x.First(webhook, id).Error
	if IsRecordNotFound(err) {
		return nil, errors.WebhookNotFound{id}
	} else if err != nil {
		return nil, err
	}
	return webhook, nil
}

type QueryWebhookOptions struct {
	CollectorID int64
	Owner       string
	RepoName    string
	EventType   string
	Sender      string
	After       int64
	Limit       int64
	Order       string
}

func QueryWebhooks(opts QueryWebhookOptions) ([]*Webhook, error) {
	db := x.Where("created > ?", opts.After)
	if opts.CollectorID > 0 {
		db = db.Where("collector_id = ?", opts.CollectorID)
	}
	if len(opts.Owner) > 0 {
		db = db.Where("owner = ?", opts.Owner)
	}
	if len(opts.RepoName) > 0 {
		db = db.Where("repo_name = ?", opts.RepoName)
	}
	if len(opts.EventType) > 0 {
		db = db.Where("event_type = ?", opts.EventType)
	}
	if len(opts.Sender) > 0 {
		db = db.Where("sender = ?", opts.Sender)
	}
	if opts.Limit > 0 {
		db = db.Limit(opts.Limit)
	}
	if len(opts.Order) > 0 {
		db = db.Order(opts.Order)
	}

	webhooks := make([]*Webhook, 0, 10)
	return webhooks, db.Find(&webhooks).Error
}

func CountWebhook() int64 {
	var count int64
	if err := x.Model(new(Webhook)).Count(&count).Error; err != nil {
		log.Printf("ERROR: CountWebhook: %v", err)
	}
	return count
}

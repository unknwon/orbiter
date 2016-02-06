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
	"fmt"
	"time"

	"github.com/Unknwon/com"

	"github.com/Unknwon/orbiter/modules/tool"
)

type CollectType int

const (
	COLLECT_TYPE_GITHUB CollectType = iota + 1
)

func (t CollectType) String() string {
	switch t {
	case COLLECT_TYPE_GITHUB:
		return "GitHub"
	}
	return com.ToStr(t)
}

// Collector represents a type of webhook collection to be stored.
type Collector struct {
	ID      int64
	Name    string `sql:"unique"`
	Type    CollectType
	Secret  string `sql:"unique"`
	Created int64
}

func (c *Collector) CreatedTime() time.Time {
	return time.Unix(0, c.Created)
}

func NewCollector(name string, tp CollectType) (*Collector, error) {
	if !x.Where("name = ?", name).First(new(Collector)).RecordNotFound() {
		return nil, ErrCollectorExists{name}
	}

	collector := &Collector{
		Name:    name,
		Type:    tp,
		Secret:  tool.NewSecretToekn(),
		Created: time.Now().UTC().UnixNano(),
	}
	if err := x.Create(collector).Error; err != nil {
		return nil, fmt.Errorf("Create: %v", err)
	}
	return collector, nil
}

func GetCollectorByID(id int64) (*Collector, error) {
	collector := new(Collector)
	err := x.First(collector, id).Error
	if IsRecordNotFound(err) {
		return nil, ErrCollectorNotFound{id, "", ""}
	} else if err != nil {
		return nil, err
	}
	return collector, nil
}

func GetCollectorBySecret(secret string) (*Collector, error) {
	collector := new(Collector)
	err := x.Where("secret = ?", secret).First(collector).Error
	if IsRecordNotFound(err) {
		return nil, ErrCollectorNotFound{0, "", secret}
	} else if err != nil {
		return nil, err
	}
	return collector, nil
}

func ListCollectors() ([]*Collector, error) {
	collectors := make([]*Collector, 0, 5)
	return collectors, x.Order("id asc").Find(&collectors).Error
}

func RegenerateCollectorSecret(id int64) error {
	return x.First(new(Collector), id).Update("secret", tool.NewSecretToekn()).Error
}

func UpdateCollector(collector *Collector) error {
	if !x.Where("name = ? AND id != ?", collector.Name, collector.ID).First(new(Collector)).RecordNotFound() {
		return ErrCollectorExists{collector.Name}
	}
	return x.Save(collector).Error
}

// TODO: delete webhook history
func DeleteCollectorByID(id int64) error {
	return x.Delete(new(Collector), id).Error
}

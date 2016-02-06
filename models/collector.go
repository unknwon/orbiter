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

type ErrCollectorExists struct {
	Name string
}

func IsErrCollectorExists(err error) bool {
	_, ok := err.(ErrCollectorExists)
	return ok
}

func (err ErrCollectorExists) Error() string {
	return fmt.Sprintf("Collector already exists: [name: %s]", err.Name)
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

func ListCollectors() ([]*Collector, error) {
	collectors := make([]*Collector, 0, 5)
	return collectors, x.Order("id asc").Find(&collectors).Error
}

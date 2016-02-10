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

package errors

import (
	"fmt"
)

type CollectorExists struct {
	Name string
}

func IsCollectorExists(err error) bool {
	_, ok := err.(CollectorExists)
	return ok
}

func (err CollectorExists) Error() string {
	return fmt.Sprintf("Collector already exists: [name: %s]", err.Name)
}

type CollectorNotFound struct {
	ID     int64
	Name   string
	Secret string
}

func IsCollectorNotFound(err error) bool {
	_, ok := err.(CollectorNotFound)
	return ok
}

func (err CollectorNotFound) Error() string {
	return fmt.Sprintf("Collector not found: [id: %d, name: %s, secret: %s]", err.ID, err.Name, err.Secret)
}

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
)

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

type ErrCollectorNotFound struct {
	ID     int64
	Name   string
	Secret string
}

func IsErrCollectorNotFound(err error) bool {
	_, ok := err.(ErrCollectorNotFound)
	return ok
}

func (err ErrCollectorNotFound) Error() string {
	return fmt.Sprintf("Collector not found: [id: %d, name: %s, secret: %s]", err.ID, err.Name, err.Secret)
}

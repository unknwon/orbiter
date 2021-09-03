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

type ApplicationExists struct {
	Name string
}

func IsApplicationExists(err error) bool {
	_, ok := err.(ApplicationExists)
	return ok
}

func (err ApplicationExists) Error() string {
	return fmt.Sprintf("Application already exists: [name: %s]", err.Name)
}

type ApplicationNotFound struct {
	ID    int64
	Name  string
	Token string
}

func IsApplicationNotFound(err error) bool {
	_, ok := err.(ApplicationNotFound)
	return ok
}

func (err ApplicationNotFound) Error() string {
	return fmt.Sprintf("Application not found: [id: %d, name: %s, token: %s]", err.ID, err.Name, err.Token)
}

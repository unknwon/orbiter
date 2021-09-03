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

type WebhookNotFound struct {
	ID int64
}

func IsWebhookNotFound(err error) bool {
	_, ok := err.(WebhookNotFound)
	return ok
}

func (err WebhookNotFound) Error() string {
	return fmt.Sprintf("Webhook not found: [id: %d]", err.ID)
}

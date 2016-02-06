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

package webhook

import (
	"encoding/json"
)

type GitHubUser struct {
	Login string `json:"login"`
}

type GitHubRepository struct {
	Name  string      `json:"name"`
	Owner *GitHubUser `json:"owner"`
}

type GitHubEvent struct {
	Repository *GitHubRepository `json:"repository"`
	Sender     *GitHubUser       `json:"sender"`
}

func ParseGitHubEvent(payload []byte) (*GitHubEvent, error) {
	event := new(GitHubEvent)
	return event, json.Unmarshal(payload, event)
}

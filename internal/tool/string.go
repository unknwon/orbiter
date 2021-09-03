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

package tool

// FirstNonEmptyString returns the value of first string
// that is not empty in the list.
// If all strings are empty, it returns an empty result.
func FirstNonEmptyString(strs ...string) string {
	for i := range strs {
		if len(strs[i]) > 0 {
			return strs[i]
		}
	}
	return ""
}

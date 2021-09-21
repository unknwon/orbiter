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

package context

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

type Context struct {
	*macaron.Context
	Flash   *session.Flash
	Session session.Store
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash) {
		ctx := &Context{
			Context: c,
			Flash:   f,
			Session: sess,
		}
		c.Map(ctx)
	}
}

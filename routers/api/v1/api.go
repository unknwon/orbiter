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

package v1

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"github.com/Unknwon/orbiter/models"
	"github.com/Unknwon/orbiter/modules/context"
)

type Context struct {
	*context.Context
	App *models.Application
}

func Contexter() macaron.Handler {
	return func(c *context.Context) {
		app, err := models.GetApplicationByToken(c.Query("token"))
		if err != nil {
			if models.IsErrApplicationExists(err) {
				c.Error(403)
			} else {
				c.Error(500, err.Error())
			}
			return
		}

		ctx := &Context{
			Context: c,
			App:     app,
		}
		c.Map(ctx)
	}
}

func RegisterRoutes(m *macaron.Macaron) {
	bind := binding.Bind
	_ = bind

	m.Group("/v1", func() {
		m.Get("/webhooks", ListWebhooks)
	}, Contexter())
}

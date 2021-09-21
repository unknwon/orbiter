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

// Orbiter is a tool for collecting and redistributing webhooks over the network.
package main

import (
	"fmt"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	log "github.com/sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"unknwon.dev/orbiter/internal/context"
	"unknwon.dev/orbiter/internal/route"
	apiv1 "unknwon.dev/orbiter/internal/route/api/v1"
	"unknwon.dev/orbiter/internal/setting"
	"unknwon.dev/orbiter/internal/templateutil"
)

func main() {
	log.Printf("Orbiter %s", setting.Version)
	m := macaron.Classic()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs:      templateutil.NewFuncMap(),
		IndentJSON: macaron.Env != macaron.PROD,
	}))
	m.Use(session.Sessioner())
	m.Use(context.Contexter())

	bindIgnErr := binding.BindIgnErr

	m.Group("", func() {
		f := flamego.New()
		f.Use(template.Templater(
			template.Options{
				FuncMaps: templateutil.NewFuncMap(),
			},
		))
		f.Get("/", route.Dashboard)
		m.Get("/", f.ServeHTTP)

		m.Group("/collectors", func() {
			m.Get("", route.Collectors)
			m.Combo("/new").Get(route.NewCollector).
				Post(bindIgnErr(route.NewCollectorForm{}), route.NewCollectorPost)
			m.Group("/:id", func() {
				m.Combo("").Get(route.EditCollector).
					Post(bindIgnErr(route.NewCollectorForm{}), route.EditCollectorPost)
				m.Post("/regenerate_token", route.RegenerateCollectorSecret)
				m.Post("/delete", route.DeleteCollector)
			})
		})

		m.Group("/applications", func() {
			m.Get("", route.Applications)
			m.Combo("/new").Get(route.NewApplication).
				Post(bindIgnErr(route.NewApplicationForm{}), route.NewApplicationPost)
			m.Group("/:id", func() {
				m.Combo("").Get(route.EditApplication).
					Post(bindIgnErr(route.NewApplicationForm{}), route.EditApplicationPost)
				m.Post("/regenerate_token", route.RegenerateApplicationSecret)
				m.Post("/delete", route.DeleteApplication)
			})
		})

		m.Group("/webhooks", func() {
			m.Get("", route.Webhooks)
			m.Get("/:id", route.ViewWebhook)
		})

		f.Get("/config", route.Config)
		m.Get("/config", f.ServeHTTP)
	}, context.BasicAuth())

	m.Post("/hook", route.Hook)

	m.Group("/api", func() {
		apiv1.RegisterRoutes(m)
	})

	listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)
	log.Println("Listening on", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, m))
}

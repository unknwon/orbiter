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
	gotmpl "html/template"
	"log"
	"net/http"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"github.com/Unknwon/orbiter/modules/context"
	"github.com/Unknwon/orbiter/modules/setting"
	"github.com/Unknwon/orbiter/modules/template"
	"github.com/Unknwon/orbiter/routers"
	apiv1 "github.com/Unknwon/orbiter/routers/api/v1"
)

const APP_VER = "0.5.1.0206"

func init() {
	setting.AppVer = APP_VER
}

func main() {
	log.Printf("Orbiter %s", APP_VER)
	m := macaron.Classic()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs:      []gotmpl.FuncMap{template.Funcs},
		IndentJSON: macaron.Env != macaron.PROD,
	}))
	m.Use(session.Sessioner())
	m.Use(context.Contexter())

	bindIgnErr := binding.BindIgnErr

	m.Group("", func() {
		m.Get("/", routers.Dashboard)

		m.Group("/collectors", func() {
			m.Get("", routers.Collectors)
			m.Combo("/new").Get(routers.NewCollector).
				Post(bindIgnErr(routers.NewCollectorForm{}), routers.NewCollectorPost)
			m.Group("/:id", func() {
				m.Combo("").Get(routers.EditCollector).
					Post(bindIgnErr(routers.NewCollectorForm{}), routers.EditCollectorPost)
				m.Post("/regenerate_token", routers.RegenerateCollectorSecret)
				m.Post("/delete", routers.DeleteCollector)
			})
		})

		m.Group("/applications", func() {
			m.Get("", routers.Applications)
			m.Combo("/new").Get(routers.NewApplication).
				Post(bindIgnErr(routers.NewApplicationForm{}), routers.NewApplicationPost)
			m.Group("/:id", func() {
				m.Combo("").Get(routers.EditApplication).
					Post(bindIgnErr(routers.NewApplicationForm{}), routers.EditApplicationPost)
				m.Post("/regenerate_token", routers.RegenerateApplicationSecret)
				m.Post("/delete", routers.DeleteApplication)
			})
		})

		m.Group("/webhooks", func() {
			m.Get("", routers.Webhooks)
			m.Get("/:id", routers.ViewWebhook)
		})

		m.Get("/config", routers.Config)
	}, context.BasicAuth())

	m.Post("/hook", routers.Hook)

	m.Group("/api", func() {
		apiv1.RegisterRoutes(m)
	})

	listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)
	log.Println("Listening on", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, m))
}

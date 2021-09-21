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

	"github.com/flamego/auth"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/flamego/template"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"

	"unknwon.dev/orbiter/internal/context"
	"unknwon.dev/orbiter/internal/route"
	apiv1 "unknwon.dev/orbiter/internal/route/api/v1"
	"unknwon.dev/orbiter/internal/setting"
	"unknwon.dev/orbiter/internal/templateutil"
)

func main() {
	if err := log.NewConsole(); err != nil {
		panic("error init logger: " + err.Error())
	}
	defer log.Stop()

	log.Info("Orbiter: %s", setting.Version)

	m := macaron.Classic()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs:      templateutil.NewFuncMap(),
		IndentJSON: macaron.Env != macaron.PROD,
	}))
	m.Use(session.Sessioner())
	m.Use(context.Contexter())

	f := flamego.New() // TODO: Use flamego.Classic()
	f.Use(template.Templater(
		template.Options{
			FuncMaps: templateutil.NewFuncMap(),
		},
	))

	f.Group("",
		func() {
			f.Get("/", route.Dashboard)

			f.Group("/collectors", func() {
				f.Get("", route.Collectors)
				f.Combo("/new").
					Get(route.NewCollector).
					Post(binding.Form(route.NewCollectorForm{}), route.NewCollectorPost)
				f.Group("/{id}", func() {
					f.Combo("").
						Get(route.EditCollector).
						Post(binding.Form(route.NewCollectorForm{}), route.EditCollectorPost)
					f.Post("/regenerate_token", route.RegenerateCollectorSecret)
					f.Post("/delete", route.DeleteCollector)
				})
			})

			f.Group("/applications", func() {
				f.Get("", route.Applications)
				f.Combo("/new").
					Get(route.NewApplication).
					Post(binding.Form(route.NewApplicationForm{}), route.NewApplicationPost)
				f.Group("/{id}", func() {
					f.Combo("").
						Get(route.EditApplication).
						Post(binding.Form(route.NewApplicationForm{}), route.EditApplicationPost)
					f.Post("/regenerate_token", route.RegenerateApplicationSecret)
					f.Post("/delete", route.DeleteApplication)
				})
			})

			f.Group("/webhooks", func() {
				f.Get("", route.Webhooks)
				f.Get("/{id}", route.ViewWebhook)
			})

			f.Get("/config", route.Config)
		},
		auth.Basic(setting.BasicAuth.User, setting.BasicAuth.Password),
	)

	m.Post("/hook", route.Hook)

	m.Group("/api", func() {
		apiv1.RegisterRoutes(m)
	})

	m.Any("*", f.ServeHTTP)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)
	log.Info("Listening on http://%s...", listenAddr)
	log.Fatal("Failed to start server: %v", http.ListenAndServe(listenAddr, m))
}

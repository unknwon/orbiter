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
)

const APP_VER = "0.1.0.0205"

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
		})

		m.Get("/config", routers.Config)
	}, context.BasicAuth())

	listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)
	log.Println("Listening on", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, m))
}

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

package route

import (
	"fmt"
	"net/http"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"unknwon.dev/orbiter/internal/db"
	"unknwon.dev/orbiter/internal/db/errors"
	"unknwon.dev/orbiter/internal/form"
)

// GET /collectors
func Collectors(w http.ResponseWriter, t template.Template, data template.Data) {
	data["Title"] = "Collectors"
	data["PageIsCollector"] = true

	collectors, err := db.ListCollectors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Collectors"] = collectors

	t.HTML(http.StatusOK, "collector/list")
}

// GET /collectors/new
func NewCollector(t template.Template, data template.Data) {
	data["Title"] = "New Collector"
	data["PageIsCollector"] = true
	t.HTML(http.StatusOK, "collector/new")
}

type NewCollectorForm struct {
	Name string `form:"name" validate:"required,max=30" name:"Collector name"`
}

// POST /collectors/new
func NewCollectorPost(c flamego.Context, t template.Template, data template.Data, errs binding.Errors, f NewCollectorForm) {
	data["Title"] = "New Collector"
	data["PageIsCollector"] = true

	if len(errs) > 0 {
		form.Validate(errs, data, f)
		t.HTML(http.StatusOK, "collector/new")
		return
	}

	collector, err := db.NewCollector(f.Name, db.CollectTypeGitHub)
	if err != nil {
		if errors.IsCollectorExists(err) {
			form.AssignForm(f, data)
			data["Error"] = "Collector name has been used."
			data["Err_Name"] = true
			t.HTML(http.StatusOK, "collector/new")
		} else {
			http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
		return
	}

	c.Redirect(fmt.Sprintf("/collectors/%d", collector.ID))
}

func parseCollectorByID(c flamego.Context) *db.Collector {
	collector, err := db.GetCollectorByID(c.ParamInt64("id"))
	if err != nil {
		if errors.IsCollectorNotFound(err) {
			http.NotFound(c.ResponseWriter(), c.Request().Request)
		} else {
			http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	return collector
}

// GET /collectors/{id}
func EditCollector(c flamego.Context, t template.Template, data template.Data) {
	data["Title"] = "Edit Collector"
	data["PageIsCollector"] = true

	data["Collector"] = parseCollectorByID(c)
	if c.ResponseWriter().Written() {
		return
	}

	t.HTML(http.StatusOK, "collector/edit")
}

// POST /collectors/{id}
func EditCollectorPost(c flamego.Context, t template.Template, data template.Data, f NewCollectorForm) {
	data["Title"] = "Edit Collector"
	data["PageIsCollector"] = true

	collector := parseCollectorByID(c)
	if c.ResponseWriter().Written() {
		return
	}
	data["Collector"] = collector

	collector.Name = f.Name
	if err := db.UpdateCollector(collector); err != nil {
		if errors.IsCollectorExists(err) {
			data["Error"] = "Collector name has been used."
			data["Err_Name"] = true
			t.HTML(http.StatusOK, "collector/edit")
		} else {
			http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
		return
	}

	c.Redirect(fmt.Sprintf("/collectors/%d", collector.ID))
}

// POST /collectors/{id}/regenerate_token
func RegenerateCollectorSecret(c flamego.Context) {
	id := c.ParamInt64("id")
	if err := db.RegenerateCollectorSecret(id); err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.Redirect(fmt.Sprintf("/collectors/%d", id))
}

// POST /collectors/{id}/delete
func DeleteCollector(c flamego.Context) {
	if err := db.DeleteCollectorByID(c.ParamInt64("id")); err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.Redirect("/collectors")
}

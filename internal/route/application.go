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

// GET /applications
func Applications(w http.ResponseWriter, t template.Template, data template.Data) {
	data["Title"] = "Applications"
	data["PageIsApplication"] = true

	apps, err := db.ListApplications()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Applications"] = apps

	t.HTML(http.StatusOK, "application/list")
}

// GET /applications/new
func NewApplication(t template.Template, data template.Data) {
	data["Title"] = "New Application"
	data["PageIsApplication"] = true
	t.HTML(http.StatusOK, "application/new")
}

type NewApplicationForm struct {
	Name string `form:"name" validate:"required,max=30" name:"Application name"`
}

// POST /applications/new
func NewApplicationPost(c flamego.Context, t template.Template, data template.Data, errs binding.Errors, f NewApplicationForm) {
	data["Title"] = "New Application"
	data["PageIsApplication"] = true

	if len(errs) > 0 {
		form.Validate(errs, data, f)
		t.HTML(http.StatusOK, "application/new")
		return
	}

	app, err := db.NewApplication(f.Name)
	if err != nil {
		if errors.IsApplicationExists(err) {
			form.AssignForm(f, data)
			data["Error"] = "Application name has been used."
			data["Err_Name"] = true
			t.HTML(http.StatusOK, "application/new")
		} else {
			http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
		return
	}

	c.Redirect(fmt.Sprintf("/applications/%d", app.ID))
}

func parseApplicationByID(c flamego.Context) *db.Application {
	app, err := db.GetApplicationByID(c.ParamInt64("id"))
	if err != nil {
		if errors.IsApplicationNotFound(err) {
			http.NotFound(c.ResponseWriter(), c.Request().Request)
		} else {
			http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	return app
}

// GET /applications/{id}
func EditApplication(c flamego.Context, t template.Template, data template.Data) {
	data["Title"] = "Edit Application"
	data["PageIsApplication"] = true

	data["Application"] = parseApplicationByID(c)
	if c.ResponseWriter().Written() {
		return
	}

	t.HTML(http.StatusOK, "application/edit")
}

// POST /applications/{id}
func EditApplicationPost(c flamego.Context, t template.Template, data template.Data, f NewApplicationForm) {
	data["Title"] = "Edit Application"
	data["PageIsApplication"] = true

	app := parseApplicationByID(c)
	if c.ResponseWriter().Written() {
		return
	}
	data["Application"] = app

	app.Name = f.Name
	if err := db.UpdateApplication(app); err != nil {
		if errors.IsApplicationExists(err) {
			data["Error"] = "Application name has been used."
			data["Err_Name"] = true
			t.HTML(http.StatusOK, "application/edit")
		} else {
			http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
		return
	}

	c.Redirect(fmt.Sprintf("/applications/%d", app.ID))
}

// POST /applications/{id}/regenerate_token
func RegenerateApplicationSecret(c flamego.Context) {
	id := c.ParamInt64("id")
	if err := db.RegenerateApplicationToken(id); err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.Redirect(fmt.Sprintf("/applications/%d", id))
}

// POST /applications/{id}/delete
func DeleteApplication(c flamego.Context) {
	if err := db.DeleteApplicationByID(c.ParamInt64("id")); err != nil {
		http.Error(c.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return
	}

	c.Redirect("/applications")
}

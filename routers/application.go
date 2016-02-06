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

package routers

import (
	"fmt"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"github.com/Unknwon/orbiter/models"
	"github.com/Unknwon/orbiter/modules/context"
	"github.com/Unknwon/orbiter/modules/form"
)

func Applications(ctx *context.Context) {
	ctx.Data["Title"] = "Applications"
	ctx.Data["PageIsApplication"] = true

	apps, err := models.ListApplications()
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Data["Applications"] = apps

	ctx.HTML(200, "application/list")
}

func NewApplication(ctx *context.Context) {
	ctx.Data["Title"] = "New Application"
	ctx.Data["PageIsApplication"] = true
	ctx.HTML(200, "application/new")
}

type NewApplicationForm struct {
	Name string `binding:"Required;MaxSize(30)" name:"Application name"`
}

func (f *NewApplicationForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return form.Validate(errs, ctx.Data, f)
}

func NewApplicationPost(ctx *context.Context, form NewApplicationForm) {
	ctx.Data["Title"] = "New Application"
	ctx.Data["PageIsApplication"] = true

	if ctx.HasError() {
		ctx.HTML(200, "application/new")
		return
	}

	app, err := models.NewApplication(form.Name)
	if err != nil {
		if models.IsErrApplicationExists(err) {
			ctx.Data["Err_Name"] = true
			ctx.RenderWithErr("Application name has been used.", "application/new", form)
		} else {
			ctx.Error(500, err.Error())
		}
		return
	}

	ctx.Redirect(fmt.Sprintf("/applications/%d", app.ID))
}

func parseApplicationByID(ctx *context.Context) *models.Application {
	app, err := models.GetApplicationByID(ctx.ParamsInt64(":id"))
	if err != nil {
		if models.IsErrApplicationNotFound(err) {
			ctx.Handle(404, "EditApplication", nil)
		} else {
			ctx.Error(500, err.Error())
		}
		return nil
	}
	ctx.Data["Application"] = app
	return app
}

func EditApplication(ctx *context.Context) {
	ctx.Data["Title"] = "Edit Application"
	ctx.Data["PageIsApplication"] = true

	parseApplicationByID(ctx)
	if ctx.Written() {
		return
	}

	ctx.HTML(200, "application/edit")
}

func EditApplicationPost(ctx *context.Context, form NewApplicationForm) {
	ctx.Data["Title"] = "Edit Application"
	ctx.Data["PageIsApplication"] = true

	app := parseApplicationByID(ctx)
	if ctx.Written() {
		return
	}

	app.Name = form.Name
	if err := models.UpdateApplication(app); err != nil {
		if models.IsErrApplicationExists(err) {
			ctx.Data["Err_Name"] = true
			ctx.RenderWithErr("Application name has been used.", "application/edit", form)
		} else {
			ctx.Error(500, err.Error())
		}
		return
	}

	ctx.Redirect(fmt.Sprintf("/applications/%d", app.ID))
}

func RegenerateApplicationSecret(ctx *context.Context) {
	if err := models.RegenerateApplicationToken(ctx.ParamsInt64(":id")); err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Redirect(fmt.Sprintf("/applications/%d", ctx.ParamsInt64(":id")))
}

func DeleteApplication(ctx *context.Context) {
	if err := models.DeleteApplicationByID(ctx.ParamsInt64(":id")); err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Redirect("/applications")
}

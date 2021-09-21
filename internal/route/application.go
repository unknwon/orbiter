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

	"unknwon.dev/orbiter/internal/context"
	"unknwon.dev/orbiter/internal/db"
	"unknwon.dev/orbiter/internal/db/errors"
)

func Applications(ctx *context.Context) {
	ctx.Data["Title"] = "Applications"
	ctx.Data["PageIsApplication"] = true

	apps, err := db.ListApplications()
	if err != nil {
		ctx.Handle(500, "ListApplications", err)
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

func NewApplicationPost(ctx *context.Context, form NewApplicationForm) {
	ctx.Data["Title"] = "New Application"
	ctx.Data["PageIsApplication"] = true

	if ctx.HasError() {
		ctx.HTML(200, "application/new")
		return
	}

	app, err := db.NewApplication(form.Name)
	if err != nil {
		if errors.IsApplicationExists(err) {
			ctx.Data["Err_Name"] = true
			ctx.RenderWithErr("Application name has been used.", "application/new", form)
		} else {
			ctx.Handle(500, "NewApplication", err)
		}
		return
	}

	ctx.Redirect(fmt.Sprintf("/applications/%d", app.ID))
}

func parseApplicationByID(ctx *context.Context) *db.Application {
	app, err := db.GetApplicationByID(ctx.ParamsInt64(":id"))
	if err != nil {
		if errors.IsApplicationNotFound(err) {
			ctx.Handle(404, "EditApplication", nil)
		} else {
			ctx.Handle(500, "GetApplicationByID", err)
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
	if err := db.UpdateApplication(app); err != nil {
		if errors.IsApplicationExists(err) {
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
	if err := db.RegenerateApplicationToken(ctx.ParamsInt64(":id")); err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Redirect(fmt.Sprintf("/applications/%d", ctx.ParamsInt64(":id")))
}

func DeleteApplication(ctx *context.Context) {
	if err := db.DeleteApplicationByID(ctx.ParamsInt64(":id")); err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Redirect("/applications")
}

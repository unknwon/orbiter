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

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"unknwon.dev/orbiter/internal/context"
	"unknwon.dev/orbiter/internal/form"
	"unknwon.dev/orbiter/internal/models"
	"unknwon.dev/orbiter/internal/models/errors"
)

func Collectors(ctx *context.Context) {
	ctx.Data["Title"] = "Collectors"
	ctx.Data["PageIsCollector"] = true

	collectors, err := models.ListCollectors()
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Data["Collectors"] = collectors

	ctx.HTML(200, "collector/list")
}

func NewCollector(ctx *context.Context) {
	ctx.Data["Title"] = "New Collector"
	ctx.Data["PageIsCollector"] = true
	ctx.HTML(200, "collector/new")
}

type NewCollectorForm struct {
	Name string `binding:"Required;MaxSize(30)" name:"Collector name"`
}

func (f *NewCollectorForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return form.Validate(errs, ctx.Data, f)
}

func NewCollectorPost(ctx *context.Context, form NewCollectorForm) {
	ctx.Data["Title"] = "New Collector"
	ctx.Data["PageIsCollector"] = true

	if ctx.HasError() {
		ctx.HTML(200, "collector/new")
		return
	}

	collector, err := models.NewCollector(form.Name, models.CollectTypeGitHub)
	if err != nil {
		if errors.IsCollectorExists(err) {
			ctx.Data["Err_Name"] = true
			ctx.RenderWithErr("Collector name has been used.", "collector/new", form)
		} else {
			ctx.Handle(500, "NewCollector", err)
		}
		return
	}

	ctx.Redirect(fmt.Sprintf("/collectors/%d", collector.ID))
}

func parseCollectorByID(ctx *context.Context) *models.Collector {
	collector, err := models.GetCollectorByID(ctx.ParamsInt64(":id"))
	if err != nil {
		if errors.IsCollectorNotFound(err) {
			ctx.Handle(404, "EditApplication", nil)
		} else {
			ctx.Handle(500, "GetCollectorByID", err)
		}
		return nil
	}
	ctx.Data["Collector"] = collector
	return collector
}

func EditCollector(ctx *context.Context) {
	ctx.Data["Title"] = "Edit Collector"
	ctx.Data["PageIsCollector"] = true

	parseCollectorByID(ctx)
	if ctx.Written() {
		return
	}

	ctx.HTML(200, "collector/edit")
}

func EditCollectorPost(ctx *context.Context, form NewCollectorForm) {
	ctx.Data["Title"] = "Edit Collector"
	ctx.Data["PageIsCollector"] = true

	collector := parseCollectorByID(ctx)
	if ctx.Written() {
		return
	}

	collector.Name = form.Name
	if err := models.UpdateCollector(collector); err != nil {
		if errors.IsCollectorExists(err) {
			ctx.Data["Err_Name"] = true
			ctx.RenderWithErr("Collector name has been used.", "collector/edit", form)
		} else {
			ctx.Error(500, err.Error())
		}
		return
	}

	ctx.Redirect(fmt.Sprintf("/collectors/%d", collector.ID))
}

func RegenerateCollectorSecret(ctx *context.Context) {
	if err := models.RegenerateCollectorSecret(ctx.ParamsInt64(":id")); err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Redirect(fmt.Sprintf("/collectors/%d", ctx.ParamsInt64(":id")))
}

func DeleteCollector(ctx *context.Context) {
	if err := models.DeleteCollectorByID(ctx.ParamsInt64(":id")); err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Redirect("/collectors")
}

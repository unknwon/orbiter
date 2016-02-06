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

package context

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"github.com/Unknwon/orbiter/modules/form"
	"github.com/Unknwon/orbiter/modules/setting"
)

type Context struct {
	*macaron.Context
	Flash   *session.Flash
	Session session.Store
}

// HasError returns true if error occurs in form validation.
func (ctx *Context) HasError() bool {
	hasErr, ok := ctx.Data["HasError"]
	if !ok {
		return false
	}
	ctx.Data["FlashTitle"] = "Form Validation"
	ctx.Flash.ErrorMsg = ctx.Data["ErrorMsg"].(string)
	ctx.Data["Flash"] = ctx.Flash
	return hasErr.(bool)
}

// RenderWithErr used for page has form validation but need to prompt error to users.
func (ctx *Context) RenderWithErr(msg string, tpl string, userForm interface{}) {
	if userForm != nil {
		form.AssignForm(userForm, ctx.Data)
	}
	ctx.Data["FlashTitle"] = "Form Validation"
	ctx.Flash.ErrorMsg = msg
	ctx.Data["Flash"] = ctx.Flash
	ctx.HTML(200, tpl)
}

// Handle handles and logs error by given status.
func (ctx *Context) Handle(status int, title string, err error) {
	if err != nil {
		log.Printf("%s: %v", title, err)
		if macaron.Env != macaron.PROD {
			ctx.Data["ErrorMsg"] = err
		}
	}

	switch status {
	case 404:
		ctx.Data["Title"] = "Page Not Found"
	case 500:
		ctx.Data["Title"] = "Internal Server Error"
	}
	ctx.HTML(status, fmt.Sprintf("status/%d", status))
}

func BasicAuthDecode(encoded string) (string, string, error) {
	s, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}

	auth := strings.SplitN(string(s), ":", 2)
	return auth[0], auth[1], nil
}

// BasicAuth preforms HTTP basic auth check.
func BasicAuth() macaron.Handler {
	return func(ctx *Context) {
		authHead := ctx.Req.Header.Get("Authorization")
		if len(authHead) == 0 {
			ctx.Resp.Header().Set("WWW-Authenticate", "Basic realm=\".\"")
			ctx.Error(401)
			return
		}

		authFields := strings.Fields(authHead)
		if len(authFields) != 2 || authFields[0] != "Basic" {
			ctx.Error(401)
			return
		}

		authUser, authPassword, err := BasicAuthDecode(authFields[1])
		if err != nil {
			ctx.Error(401)
			return
		}
		if authUser != setting.BasicAuth.User ||
			authPassword != setting.BasicAuth.Password {
			ctx.Error(401)
			return
		}
	}
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash) {
		ctx := &Context{
			Context: c,
			Flash:   f,
			Session: sess,
		}
		c.Map(ctx)
	}
}

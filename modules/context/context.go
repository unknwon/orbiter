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
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/Unknwon/orbiter/modules/setting"
)

type Context struct {
	*macaron.Context
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
	return func(c *macaron.Context) {
		ctx := &Context{
			Context: c,
		}
		c.Map(ctx)
	}
}

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

package setting

import (
	"gopkg.in/ini.v1"
	log "unknwon.dev/clog/v2"

	"unknwon.dev/orbiter/internal/tool"
)

var (
	Version = "dev"

	HTTPPort int

	BasicAuth struct {
		User     string
		Password string
	}

	Database struct {
		Host     string
		Name     string
		User     string
		Password string
	}

	Cfg *ini.File
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatal("Fail to load configuration: %v", err)
	}
	if tool.IsFile("custom/app.ini") {
		if err = Cfg.Append("custom/app.ini"); err != nil {
			log.Fatal("Fail to load custom configuration: %v", err)
		}
	}
	Cfg.NameMapper = ini.SnackCase

	HTTPPort = Cfg.Section("").Key("HTTP_PORT").MustInt(8085)
	if err = Cfg.Section("basic_auth").MapTo(&BasicAuth); err != nil {
		log.Fatal("Fail to map section 'basic_auth': %v", err)
	} else if err = Cfg.Section("database").MapTo(&Database); err != nil {
		log.Fatal("Fail to map section 'database': %v", err)
	}
}

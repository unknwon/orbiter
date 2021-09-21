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
	log "github.com/sirupsen/logrus"
	"github.com/unknwon/com"
	"gopkg.in/ini.v1"
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
		log.Fatalf("Fail to load configuration: %s", err)
	}
	if com.IsFile("custom/app.ini") {
		if err = Cfg.Append("custom/app.ini"); err != nil {
			log.Fatalf("Fail to load custom configuration: %s", err)
		}
	}
	Cfg.NameMapper = ini.AllCapsUnderscore

	HTTPPort = Cfg.Section("").Key("HTTP_PORT").MustInt(8085)
	if err = Cfg.Section("basic_auth").MapTo(&BasicAuth); err != nil {
		log.Fatalf("Fail to map section 'basic_auth': %s", err)
	} else if err = Cfg.Section("database").MapTo(&Database); err != nil {
		log.Fatalf("Fail to map section 'database': %s", err)
	}
}

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

package models

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/Unknwon/orbiter/modules/setting"
)

var x *gorm.DB

func init() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
		setting.Database.User, setting.Database.Password, setting.Database.Host, setting.Database.Name))
	if err != nil {
		log.Fatalf("Fail to open database: %s", err)
	}
	x = &db

	x.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(new(Collector))
}

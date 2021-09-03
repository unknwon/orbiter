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
	"time"

	"unknwon.dev/orbiter/models/errors"
	"unknwon.dev/orbiter/modules/tool"
)

// Application represents a consumer application that calls APIs.
type Application struct {
	ID      int64
	Name    string `sql:"unique"`
	Token   string `sql:"unique"`
	Created int64
}

func (app *Application) CreatedTime() time.Time {
	return time.Unix(0, app.Created)
}

func NewApplication(name string) (*Application, error) {
	if !x.Where("name = ?", name).First(new(Application)).RecordNotFound() {
		return nil, errors.ApplicationExists{name}
	}

	app := &Application{
		Name:    name,
		Token:   tool.NewSecretToekn(),
		Created: time.Now().UTC().UnixNano(),
	}
	if err := x.Create(app).Error; err != nil {
		return nil, fmt.Errorf("Create: %v", err)
	}
	return app, nil
}

func GetApplicationByID(id int64) (*Application, error) {
	app := new(Application)
	err := x.First(app, id).Error
	if IsRecordNotFound(err) {
		return nil, errors.ApplicationNotFound{id, "", ""}
	} else if err != nil {
		return nil, err
	}
	return app, nil
}

func GetApplicationByToken(token string) (*Application, error) {
	app := new(Application)
	err := x.Where("token = ?", token).First(app).Error
	if IsRecordNotFound(err) {
		return nil, errors.ApplicationNotFound{0, "", token}
	} else if err != nil {
		return nil, err
	}
	return app, nil
}

func ListApplications() ([]*Application, error) {
	apps := make([]*Application, 0, 5)
	return apps, x.Order("id asc").Find(&apps).Error
}

func RegenerateApplicationToken(id int64) error {
	return x.First(new(Application), id).Update("token", tool.NewSecretToekn()).Error
}

func UpdateApplication(app *Application) error {
	if !x.Where("name = ? AND id != ?", app.Name, app.ID).First(new(Application)).RecordNotFound() {
		return errors.ApplicationExists{app.Name}
	}
	return x.Save(app).Error
}

func DeleteApplicationByID(id int64) error {
	return x.Delete(new(Application), id).Error
}

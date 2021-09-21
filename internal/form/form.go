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

package form

import (
	"reflect"

	"github.com/flamego/binding"
	"github.com/flamego/template"
	"github.com/flamego/validator"
)

// AssignForm assign form values back to the template data.
func AssignForm(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = field.Name
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func Validate(errs binding.Errors, data template.Data, f interface{}) {
	AssignForm(f, data)

	var validationErrs validator.ValidationErrors
	for _, err := range errs {
		if err.Category == binding.ErrorCategoryValidation {
			validationErrs = err.Err.(validator.ValidationErrors)
			break
		}
	}
	if len(validationErrs) == 0 {
		return
	}

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		}

		if validationErrs[0].StructField() == field.Name {
			data["Err_"+field.Name] = true

			name := field.Tag.Get("name")

			switch validationErrs[0].Tag() {
			case "max":
				data["Error"] = name + " cannot contain more than " + validationErrs[0].Param() + " characters."
				// 	case binding.ERR_REQUIRED:
				// 		data["ErrorMsg"] = name + " cannot be empty."
				// 	case binding.ERR_ALPHA_DASH:
				// 		data["ErrorMsg"] = name + " must be valid alpha or numeric or dash(-_) characters."
				// 	case binding.ERR_ALPHA_DASH_DOT:
				// 		data["ErrorMsg"] = name + " must be valid alpha or numeric or dash(-_) or dot characters."
				// 	case binding.ERR_SIZE:
				// 		data["ErrorMsg"] = name + " must be size " + GetSize(field)
				// 	case binding.ERR_MIN_SIZE:
				// 		data["ErrorMsg"] = name + " must contain at least " + GetMinSize(field) + " characters."
				// 	case binding.ERR_EMAIL:
				// 		data["ErrorMsg"] = name + " is not a valid email address."
				// 	case binding.ERR_URL:
				// 		data["ErrorMsg"] = name + " is not a valid URL."
				// 	case binding.ERR_INCLUDE:
				// 		data["ErrorMsg"] = name + " must contain substring '" + GetInclude(field) + "'."
			default:
				data["Error"] = validationErrs[0].Error()
			}
			break
		}
	}
}
